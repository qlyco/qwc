const cacheName = "qwc-site"
const startPage = new URL("./", self.location).href

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(cacheName).then((cache) => cache.add(startPage)))
  self.skipWaiting()
})

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim())
})

self.addEventListener("fetch", (event) => {
  const request = event.request

  if (request.method !== "GET" || new URL(request.url).origin !== self.location.origin) {
    return
  }

  event.respondWith(fromNetwork(request))
})

async function fromNetwork(request) {
  const cache = await caches.open(cacheName)

  try {
    const response = await fetch(request)

    if (response.ok) {
      cache.put(request, response.clone())
    }

    return response
  } catch {
    const cached = await cache.match(request)

    if (cached) {
      return cached
    }

    if (request.mode === "navigate") {
      return cache.match(startPage)
    }

    return Response.error()
  }
}
