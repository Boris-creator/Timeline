import './style.css'
import {fabric} from 'fabric'

type Diapason = {
    from: number
    to: number
}
type Event = {
    name: string
    description: string
    id: number
    dateFrom: Date
    dateTo: Date
    duration: number
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

const prepareEvent: (ev: Event) => Event = (ev: Event) => ({
    ...ev,
    dateFrom: new Date(ev.dateFrom),
    dateTo: new Date(ev.dateTo),
    duration: diffInDays(new Date(ev.dateFrom), new Date(ev.dateTo)),
})

/** Dates utils */
const addDays = (date: Date, days: number) => {
    const newDate = new Date(date)
    newDate.setDate(date.getDate() + days)
    return newDate
}
const diffInDays = (date1: Date, date2: Date) => (+date1 - +date2) / (1000 * 60 * 60 * 24)

/** Diapason utils */
const isIntersecting = (d1: Diapason, d2: Diapason) => Math.max(d1.from, d2.from) <= Math.min(d1.to, d2.to)
const isContainingIn = (d1: Diapason, d2: Diapason) => d1.from >= d2.from && d1.to <= d2.to

/** Constants */
const TIMELINE_LENGTH = 7
const DAYS_TOTAL = 180

const startDate = addDays(new Date(), -DAYS_TOTAL)

class Fetcher {
    private exploredArea = new Set<Diapason>()
    private events: Array<Event> = []

    public fetchEvents = debounce(async (diapason: Diapason) => {
        if (this.checkIfExplored(diapason)) {
            return []
        }
        try {
            const response = await request('/events/search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify({
                    dateFrom: addDays(startDate, diapason.from),
                    dateTo: addDays(startDate, diapason.to),
                })
            })
            const events = await response.json()

            this.addToExplored(diapason)
            return events
        } catch {
        }
    }, 500)

    public addEvents(eventsData: Array<Event>) {
        const newEvents: Array<Event> = []
        eventsData.forEach(ev => {
            if (!this.events.some(({id}) => id === ev.id)) {
                this.events.push(ev)
                newEvents.push(ev)
            }
        })
        return newEvents
    }

    private checkIfExplored(diapason: Diapason) {
        return [...this.exploredArea].some(exploredDiapason => isContainingIn(diapason, exploredDiapason))
    }

    private addToExplored(diapason: Diapason) {
        if (this.checkIfExplored(diapason)) {
            return
        }
        const intersecting = [...this.exploredArea].filter(exploredDiapason => isIntersecting(exploredDiapason, diapason))
        intersecting.forEach(d => {
            this.exploredArea.delete(d)
        })
        this.exploredArea.add({
            from: Math.min(...intersecting.map(({from}) => from), diapason.from),
            to: Math.max(...intersecting.map(({to}) => to), diapason.to)
        })
    }
}

class Renderer {
    private readonly CANVAS_WIDTH: number
    private readonly CANVAS_HEIGHT: number
    private canvas: fabric.Canvas

    private groups: Record<string, fabric.Group> = {}

    constructor() {
        const windowSize = Math.min(window.innerWidth, window.innerHeight)
        this.CANVAS_HEIGHT = windowSize
        this.CANVAS_WIDTH = windowSize
        const canvasElement = document.createElement('canvas')
        canvasElement.width = this.CANVAS_WIDTH
        canvasElement.height = this.CANVAS_HEIGHT
        document.body.append(canvasElement)
        this.canvas = new fabric.Canvas(canvasElement)

        this.drawTimeLine()
        this.initHandlers()
    }

    public renderEvent(event: Event, priority = 0) {
        const group = this.groups.timeline
        const dayStart = diffInDays(event.dateFrom, startDate)
        const daysLong = diffInDays(event.dateTo, event.dateFrom) + 1
        const eventRectWidth = 50
        const eventRect = new fabric.Rect({
            left: 10 + (eventRectWidth * 1.1) * priority,
            top: this.getOrdinateByDay(dayStart),
            width: eventRectWidth,
            height: daysLong * this.dayInterval,
            fill: 'rgba(50, 205, 50, .3)',
            padding: 10
        })
        group.add(eventRect)
        group.canvas?.add(eventRect)
        return eventRect
    }

    public renderEvents(events: Array<Event>) {
        events.forEach((event, i) => {
            const priority = events.filter(({duration, dateFrom, dateTo}, evIdx) =>
                duration >= event.duration
                && i !== evIdx
                && isIntersecting(
                    {from: +dateFrom, to: +dateTo},
                    {from: +event.dateFrom, to: +event.dateTo}
                )).length
            this.renderEvent(event, priority)
        })
    }

    public emitDiapason() {
        window.dispatchEvent(new CustomEvent<Diapason>('timeScroll', {detail: this.getDaysByPosition()}))
    }

    private initHandlers() {
        const {CANVAS_HEIGHT, canvas, groups} = this
        const that = this
        canvas.on('mouse:wheel', function (opt) {
            opt.e.preventDefault()
            opt.e.stopPropagation()
            const delta = opt.e.deltaY;
            if (opt.e.ctrlKey) {
                const coefficient = 0.999 ** delta
                const scaling = groups.timeline.getObjectScaling();
                const newZoom = {
                    scaleY: Math.min(4, Math.max(.25, scaling.scaleY * coefficient))
                }
                const scaleCoefficient = newZoom.scaleY / scaling.scaleY
                const center = CANVAS_HEIGHT / 2
                const lag = center - groups.timeline.getCenterPoint().y
                groups.timeline.set({
                    scaleY: newZoom.scaleY,
                    top: center - lag * scaleCoefficient,
                    originY: 'center'
                })
                groups.daysMarkup.visible = newZoom.scaleY > 1
            } else {
                groups.timeline.set({
                    top: (groups.timeline.top ?? 0) - delta / 2 //TODO infinite scroll
                })
            }
            that.emitDiapason()
            canvas.renderAll()
        })
    }

    private getOrdinateByDay(day: number) {
        const zeroPoint = -this.CANVAS_HEIGHT * (TIMELINE_LENGTH / 2)
        return zeroPoint + this.dayInterval * day
    }

    private get dayInterval() {
        return this.CANVAS_HEIGHT * TIMELINE_LENGTH / DAYS_TOTAL
    }

    private drawTimeLine() {
        const {CANVAS_WIDTH, CANVAS_HEIGHT, canvas} = this
        const linePoints = [
            CANVAS_WIDTH / 2,
            -CANVAS_HEIGHT * ((TIMELINE_LENGTH - 1) / 2),
            CANVAS_WIDTH / 2,
            CANVAS_HEIGHT * ((TIMELINE_LENGTH + 1) / 2)
        ]

        const timeLine = new fabric.Line(linePoints, {
            stroke: 'black',
        })
        const days: Array<fabric.Line> = []
        const weeks: Array<fabric.Line> = []
        for (let i = 1; i < DAYS_TOTAL; i++) {
            const ordinate = (linePoints[3] - linePoints[1]) / DAYS_TOTAL * i + linePoints[1]
            const coords = [
                CANVAS_WIDTH * .25,
                ordinate,
                CANVAS_WIDTH * .75,
                ordinate
            ]
            const mark = new fabric.Line(coords, {
                stroke: 'black',
            })
            days.push(mark)

            const date = addDays(startDate, i)

            if (date.getDay() === 0) {
                const weekMark = new fabric.Line([
                    0,
                    coords[1],
                    CANVAS_WIDTH,
                    coords[3]
                ], {
                    stroke: 'red',
                })
                weeks.push(weekMark)
            }
        }
        const daysMarkup = new fabric.Group(days)
        const markup = new fabric.Group([daysMarkup, ...weeks])
        const group = new fabric.Group([timeLine, markup], {
            selectable: false
        })
        canvas.add(group)

        this.groups.daysMarkup = daysMarkup
        this.groups.timeline = group
    }

    private getDaysByPosition() {
        const scale = this.groups.timeline.scaleY || 1
        const visibleTop = -this.groups.timeline.getBoundingRect(true, true).top
        const dayLong = (TIMELINE_LENGTH * this.CANVAS_HEIGHT / DAYS_TOTAL) * scale
        const diapason: Diapason = {
            from: Math.round(visibleTop / dayLong),
            to: Math.round((visibleTop + this.CANVAS_HEIGHT) / dayLong)
        }
        return diapason
    }
}

function main() {
    const renderer = new Renderer()
    const fetcher = new Fetcher()

    window.addEventListener('timeScroll', (e) => {
        const ev = e as CustomEvent<Diapason>
        fetcher.fetchEvents(ev.detail)
            .then(debounced => debounced)
            .then((data: Array<Event>) => {
                const newEvents = fetcher.addEvents(data.map(prepareEvent))
                renderer.renderEvents(newEvents)
            })
    })

    renderer.emitDiapason()
}

main()
