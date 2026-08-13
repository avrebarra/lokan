#!/usr/bin/env node
// Mock API server for lokan web dev — in-memory, zero deps.
// Implements the frozen contract in docs/api.md so the Vite dev UI is
// fully interactive without building the Go binary. Run via `runtask web dev`
// (Vite proxies /api here). Port: MOCK_API_PORT (default 8787).

import http from 'node:http'

const PORT = Number(process.env.MOCK_API_PORT ?? 8787)

// ─── in-memory store, seeded with a small airline demo board ───────────────

let nextCounter = 100
const tasks = new Map()

// configurable lanes, mirroring the real config.json defaults
let statuses = [
  { id: 'backlog', archived: false },
  { id: 'todo', archived: false },
  { id: 'in-progress', archived: false },
  { id: 'done', archived: true },
  { id: 'cancelled', archived: true },
]

function seed(partials) {
  for (const p of partials) {
    const id = `task-${nextCounter++}`
    tasks.set(id, {
      id,
      title: p.title,
      type: p.type ?? 'task',
      status: p.status ?? 'todo',
      priority: p.priority ?? 'medium',
      parent: p.parent ?? '',
      related: p.related ?? [],
      docs: p.docs ?? [],
      tags: p.tags ?? [],
      created: p.created ?? '2026-08-13',
      updated: p.updated ?? '2026-08-13',
      body: p.body ?? '',
      filePath: `/mock/tasks/${id}.md`,
    })
  }
}

seed([
  {
    title: 'Launch SFO–NRT Route',
    type: 'epic',
    status: 'in-progress',
    priority: 'critical',
    tags: ['routes', 'international'],
  },
  {
    title: 'Submit route approval to FAA and JCAB',
    status: 'done',
    priority: 'critical',
    parent: 'task-100',
    tags: ['compliance'],
  },
  {
    title: 'Assign Boeing 787-9 fleet for long-haul',
    status: 'in-progress',
    priority: 'high',
    parent: 'task-100',
    tags: ['fleet'],
  },
  {
    title: 'Hire and train Japan-route cabin crew',
    status: 'todo',
    priority: 'high',
    parent: 'task-100',
    tags: ['crew'],
  },
  {
    title: 'Launch marketing campaign for SFO–NRT',
    status: 'backlog',
    priority: 'medium',
    parent: 'task-100',
    tags: ['marketing'],
  },
  {
    title: 'Upgrade Passenger Experience',
    type: 'epic',
    status: 'in-progress',
    priority: 'high',
    tags: ['passenger'],
  },
  {
    title: 'Redesign business class seating',
    status: 'done',
    priority: 'high',
    parent: 'task-105',
    tags: ['cabin'],
  },
  {
    title: 'Roll out in-flight Wi-Fi on all A320s',
    status: 'in-progress',
    priority: 'high',
    parent: 'task-105',
    tags: ['connectivity'],
  },
  {
    title: 'Introduce premium meal service',
    status: 'todo',
    priority: 'medium',
    parent: 'task-105',
    tags: ['catering'],
  },
  {
    title: 'Check-in kiosk freezes at bag-drop screen',
    type: 'bug',
    status: 'in-progress',
    priority: 'critical',
    tags: ['kiosk'],
  },
])

// ─── helpers ────────────────────────────────────────────────────────────────

function summary(t) {
  const { body, ...rest } = t
  return rest
}

function respond(res, status, payload) {
  res.writeHead(status, { 'Content-Type': 'application/json' })
  res.end(JSON.stringify(payload))
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let raw = ''
    req.on('data', (c) => (raw += c))
    req.on('end', () => {
      try {
        resolve(raw ? JSON.parse(raw) : {})
      } catch (err) {
        reject(err)
      }
    })
  })
}

// ─── routing ────────────────────────────────────────────────────────────────

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://localhost:${PORT}`)
  const path = url.pathname

  try {
    // tasks list
    if (req.method === 'GET' && path === '/api/tasks') {
      const all = [...tasks.values()]
      respond(res, 200, { tasks: all.map(summary), statuses, root: '/mock' })
      return
    }

    // single task
    if (req.method === 'GET' && path.startsWith('/api/task/')) {
      const id = path.slice('/api/task/'.length)
      const task = tasks.get(id)
      if (!task) return respond(res, 404, { error: `task not found: ${id}` })
      respond(res, 200, { task })
      return
    }

    // create task
    if (req.method === 'POST' && path === '/api/create') {
      const body = await readBody(req)
      if (!body.title) return respond(res, 400, { error: 'Missing title' })
      const id = `task-${nextCounter++}`
      const task = {
        id,
        title: body.title,
        type: body.type ?? 'task',
        status: 'todo',
        priority: body.priority ?? 'medium',
        parent: body.parent ?? '',
        related: [],
        docs: [],
        tags: [],
        created: '2026-08-13',
        updated: '2026-08-13',
        body: `# ${body.title}\n\n## Notes\n\n\n## Work Log\n\n`,
        filePath: `/mock/tasks/${id}.md`,
      }
      tasks.set(id, task)
      respond(res, 200, { task })
      return
    }

    // update task field
    if (req.method === 'POST' && path === '/api/update') {
      const body = await readBody(req)
      const task = tasks.get(body.id)
      if (!task) return respond(res, 404, { error: `task not found: ${body.id}` })
      if (!['status', 'priority', 'title'].includes(body.field)) {
        return respond(res, 400, { error: `Unknown field: ${body.field}` })
      }
      task[body.field] = body.value
      task.updated = '2026-08-13'
      respond(res, 200, { task })
      return
    }

    // replace the lane set (mock: renames/removals move tasks to first lane)
    if (req.method === 'POST' && path === '/api/config/statuses') {
      const body = await readBody(req)
      const next = body.statuses ?? []
      if (next.length === 0) return respond(res, 400, { error: 'At least one status is required' })
      const ids = new Set(next.map((s) => s.id))
      if (ids.size !== next.length || [...ids].some((id) => !id)) {
        return respond(res, 400, { error: 'Statuses must be non-empty and unique' })
      }
      let moved = 0
      const known = new Set(statuses.map((s) => s.id))
      for (const t of tasks.values()) {
        if (!known.has(t.status)) continue
        if (!ids.has(t.status)) {
          t.status = next[0].id
          moved++
        }
      }
      statuses = next
      respond(res, 200, { statuses, moved })
      return
    }

    // bulk delete: archived lanes or the whole board
    if (req.method === 'POST' && path === '/api/clear') {
      const body = await readBody(req)
      const archivedIds = new Set(statuses.filter((s) => s.archived).map((s) => s.id))
      let deleted = 0
      for (const [id, t] of tasks) {
        if (body.scope === 'all' || (body.scope === 'archived' && archivedIds.has(t.status))) {
          tasks.delete(id)
          deleted++
        }
      }
      respond(res, 200, { deleted })
      return
    }

    // reseed demo data
    if (req.method === 'POST' && path === '/api/seed') {
      tasks.clear()
      nextCounter = 100
      seed([
        {
          title: 'Launch SFO–NRT Route',
          type: 'epic',
          status: 'in-progress',
          priority: 'critical',
        },
        {
          title: 'Submit route approval',
          status: 'done',
          priority: 'critical',
          parent: 'task-100',
        },
        {
          title: 'Assign Boeing 787-9 fleet',
          status: 'in-progress',
          priority: 'high',
          parent: 'task-100',
        },
      ])
      respond(res, 200, { created: tasks.size })
      return
    }

    respond(res, 404, { error: 'Not Found' })
  } catch (err) {
    respond(res, 500, { error: err.message })
  }
})

server.listen(PORT, () => {
  console.log(`mock lokan api — http://localhost:${PORT}`)
})
