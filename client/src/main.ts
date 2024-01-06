import './style.css'
import {fabric} from 'fabric'

type Diapason = { from: number, to: number }
type Event = {
    name: string,
    description: string,
    id: number,
}

/** Common utils */
function debounce<F extends (params: any) => any>(callback: F, ms: number) {
    let timeout: number
    return function (...args: Parameters<F>): Promise<ReturnType<F>> {
        clearTimeout(timeout);
        return new Promise(res => {
            timeout = setTimeout(() => {
                res(callback.apply(null, args))
            }, ms)
        })
    };
}

const request = (input: string, init?: RequestInit) => {
    const baseUrl = 'http://localhost:3000/api'
    return new Promise<Response>((resolve, reject) => {
        fetch(`${baseUrl}${input}`, init).then(response => {
            if (response.ok) {
                resolve(response)
            } else {
                reject(response)
            }
        })
    })
}

const addDays = (date: Date, days: number) => {
    const newDate = new Date(date)
    newDate.setDate(date.getDate() + days)
    return newDate
}

/** Diapason utils */
const isIntersecting = (d1: Diapason, d2: Diapason) => Math.max(d1.from, d2.from) <= Math.min(d1.to, d2.to)
const isContainingIn = (d1: Diapason, d2: Diapason) => d1.from >= d2.from && d1.to <= d2.to

/** Constants */
const timeLineLength = 7
const daysCount = 180

const startDate = addDays(new Date(), -daysCount)
const exploredArea = new Set<Diapason>()
const events: Array<Event> = []

const checkIfExplored = (diapason: Diapason) => [...exploredArea].some(exploredDiapason => isContainingIn(diapason, exploredDiapason))
const addToExplored = (diapason: Diapason) => {
    if (checkIfExplored(diapason)) {
        return
    }
    const intersecting = [...exploredArea].filter(exploredDiapason => isIntersecting(exploredDiapason, diapason))
    intersecting.forEach(d => {
        exploredArea.delete(d)
    })
    exploredArea.add({
        from: Math.min(...intersecting.map(({from}) => from), diapason.from),
        to: Math.max(...intersecting.map(({to}) => to), diapason.to)
    })
}

const addEvents = (eventsData: Array<Event>) => {
    eventsData.forEach(ev => {
        if (!events.some(({id}) => id === ev.id)) {
            events.push(ev)
        }
    })
}

const fetchEvents = debounce((diapason: Diapason) => {
    return request('/events/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify({
            dateFrom: addDays(startDate, diapason.from),
            dateTo: addDays(startDate, diapason.to),
        })
    })
}, 500)

function searchByDiapason(diapason: Diapason) {
    if (checkIfExplored(diapason)) {
        return
    }
    fetchEvents(diapason)
        .then(debounced => debounced)
        .then(response => response.json())
        .then((data: Array<Event>) => {
            addEvents(data)
            addToExplored(diapason)
        })
}

function searchByPosition(group: fabric.Group, totalHeight: number) {
    const scale = group.scaleY || 1
    const visibleTop = -group.getBoundingRect(true, true).top
    const dayLong = (timeLineLength * totalHeight / daysCount) * scale
    searchByDiapason({
        from: Math.round(visibleTop / dayLong),
        to: Math.round((visibleTop + totalHeight) / dayLong)
    })
}

function main() {
    const windowSize = Math.min(window.innerWidth, window.innerHeight)
    const CANVAS_WIDTH = windowSize
    const CANVAS_HEIGHT = windowSize
    const canvasElement = document.createElement('canvas')
    canvasElement.width = CANVAS_WIDTH
    canvasElement.height = CANVAS_HEIGHT
    document.body.append(canvasElement)
    const canvas = new fabric.Canvas(canvasElement)

    const linePoints = [
        CANVAS_WIDTH / 2,
        -CANVAS_HEIGHT * ((timeLineLength - 1) / 2),
        CANVAS_WIDTH / 2,
        CANVAS_HEIGHT * ((timeLineLength + 1) / 2)
    ]

    const timeLine = new fabric.Line(linePoints, {
        stroke: 'black',
    })
    const days: Array<fabric.Line> = []
    const weeks: Array<fabric.Line> = []
    for (let i = 1; i < daysCount; i++) {
        const ordinate = (linePoints[3] - linePoints[1]) / daysCount * i + linePoints[1]
        const coords = [
            CANVAS_WIDTH * .25,
            ordinate,
            CANVAS_WIDTH * .75,
            ordinate
        ]
        const test = new fabric.Line(coords, {
            stroke: 'black',
        })
        days.push(test)

        if (i % 7 === 0) {
            const test = new fabric.Line([
                0,
                coords[1],
                CANVAS_WIDTH,
                coords[3]
            ], {
                stroke: 'red',
            })
            weeks.push(test)
        }
    }
    const daysMarkup = new fabric.Group(days)
    const markup = new fabric.Group([daysMarkup, ...weeks])
    const group = new fabric.Group([timeLine, markup], {
        selectable: false
    })
    canvas.add(group)

    searchByPosition(group, CANVAS_HEIGHT)

    canvas.on('mouse:wheel', function (opt) {
        opt.e.preventDefault()
        opt.e.stopPropagation()
        const delta = opt.e.deltaY;
        if (opt.e.ctrlKey) {
            const coefficient = 0.999 ** delta
            const scaling = group.getObjectScaling();
            const newZoom = {
                scaleY: Math.min(4, Math.max(.25, scaling.scaleY * coefficient))
            }
            const scaleCoefficient = newZoom.scaleY / scaling.scaleY
            const center = CANVAS_HEIGHT / 2
            const lag = center - group.getCenterPoint().y
            group.set({
                scaleY: newZoom.scaleY,
                top: center - lag * scaleCoefficient,
                originY: 'center'
            })
            daysMarkup.visible = newZoom.scaleY > 1
        } else {
            group.set({
                top: (group.top ?? 0) - delta / 2 //TODO infinite scroll
            })
        }
        searchByPosition(group, CANVAS_HEIGHT)
        canvas.renderAll()
    })
}

main()
