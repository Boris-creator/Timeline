function test () {
  fetch('/api/events', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({
      name: "test",
      description: "test descr",
      dateFrom: new Date(),
      dateTo: new Date(),
      precisionFrom: "day",
      precisionTo: "hour",
      eventRoles: [{
        name: "witness",
        description: "Stood near and saw it all",
        eventParticipants: [{
          dateFrom: new Date(),
        }]
      }],
    })
  })
}

function testRegister() {
  fetch('/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({
      login: "root",
      password: "qwerty",
    })
  })
}

function testLogin() {
  return fetch('/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({
      login: "root",
      password: "qwerty",
    })
  })
}

testLogin().then(test)