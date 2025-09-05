# Database Intelligence UI

A modern React-based web interface for the Enterprise Database Intelligence System.

## Architecture

This UI follows **Clean Architecture** principles with:

- **Feature-based folder structure** for scalability
- **Custom hooks** for business logic separation
- **Component composition** patterns
- **Type-safe API layer** with TypeScript
- **State management** with Zustand
- **Server state** management with TanStack Query

## Technology Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **Zustand** for client state management
- **TanStack Query** for server state management
- **React Router** for navigation
- **Recharts** for data visualization
- **Lucide React** for icons

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The UI will be available at `http://localhost:3000` and will proxy API requests to `http://localhost:8080`.

### Building

```bash
npm run build
```

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # Basic UI components (Button, Input, etc.)
│   └── layout/         # Layout components (Header, Sidebar, etc.)
├── pages/              # Page components
├── hooks/              # Custom React hooks
├── services/           # API service layer
├── stores/             # Zustand stores
├── types/              # TypeScript type definitions
└── utils/              # Utility functions
```

## Key Features

- **Responsive Design**: Works on desktop, tablet, and mobile
- **Real-time Updates**: Uses React Query for efficient data fetching
- **Type Safety**: Full TypeScript coverage
- **Error Handling**: Comprehensive error boundaries and user feedback
- **Performance**: Optimized with code splitting and lazy loading
- **Accessibility**: WCAG compliant components

## API Integration

The UI communicates with the Go backend through a well-defined REST API:

- **Connections**: Manage database connections
- **Analysis**: Perform and view database analysis
- **Security**: Security scanning and vulnerability assessment
- **Optimization**: Query optimization suggestions
- **Monitoring**: Real-time monitoring and alerts
- **Lineage**: Data lineage tracking and visualization
