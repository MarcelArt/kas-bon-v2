/// <reference types="vite/client" />

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      readonly API_URL?: string
      readonly NODE_ENV: "development" | "production" | "test"
    }
  }
}

export {}
