// Preload script. Exposes a minimal, safe API to the renderer process.
const { contextBridge } = require('electron')

contextBridge.exposeInMainWorld('rentacar', {
  isDesktop: true,
  platform: process.platform,
  version: process.versions.electron,
})
