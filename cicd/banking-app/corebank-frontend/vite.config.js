import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  server: {
    proxy: {
      // Proxy all API requests to the Go backend
      '/api': {
        target: 'http://localhost:8080', // Your Go API base URL
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        secure: false,
        configure: (proxy) => {
          proxy.on('error', (err) => {
            console.error('Proxy error:', err)
          })
          proxy.on('proxyReq', (proxyReq) => {
            console.log('Proxy request to:', proxyReq.path)
          })
        }
      },
      // Optionally proxy WebSocket connections if needed
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true
      }
    }
  },
  plugins: [
    react(),
    tailwindcss({
      config: {
        content: [
          "./index.html",
          "./src/**/*.{js,ts,jsx,tsx}",
        ],
        theme: {
          extend: {
            colors: {
              primary: {
                DEFAULT: 'oklch(59.59% 0.24 255.09)',
                dark: 'oklch(47.62% 0.24 255.85)'
              },
              secondary: {
                DEFAULT: 'oklch(70% 0.2 180)'
              }
            },
            fontFamily: {
              sans: ['Inter', 'sans-serif']
            }
          }
        },
        plugins: [
          // require('@tailwindcss/forms'),
          // require('@tailwindcss/typography')
        ]
      }
    })
  ],
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom'
    ]
  },
  build: {
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
          tailwind: ['tailwindcss']
        }
      }
    }
  }
})

// import { defineConfig } from 'vite'
// import react from '@vitejs/plugin-react'
// import tailwindcss from '@tailwindcss/vite'

// export default defineConfig({
//   server: {
//     proxy: {
//       // Proxy API requests to the Go backend
//       '/api': {
//         target: 'http://localhost:8080', // Go API URL
//         changeOrigin: true,
//         rewrite: (path) => path.replace(/^\/api/, ''), // Optional: if your Go API has a common prefix like `/api`
//       },
//     },
//   },
//   plugins: [
//     react(),
//     tailwindcss({
//       config: {
//         content: [
//           "./index.html",
//           "./src/**/*.{js,ts,jsx,tsx}",
//         ],
//         theme: {
//           extend: {
//             colors: {
//               primary: {
//                 DEFAULT: 'oklch(59.59% 0.24 255.09)',
//                 dark: 'oklch(47.62% 0.24 255.85)'
//               }
//             }
//           }
//         }
//       }
//     })
//   ]
// })