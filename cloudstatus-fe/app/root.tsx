import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  Link,
  useLocation
} from "react-router";
import { Navbar, Nav, Container } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";

import { useState, useEffect } from "react";

import type { Route } from "./+types/root";
import "./app.css";

export const links: Route.LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
];

export function Layout({ children, title }: { children: React.ReactNode; title: string }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
        <title>{title}</title>
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

interface Node {
  id: string;
  label: string;
  location: string;
  reset_day: number;
}

interface ApiResponse {
  title: string;
  nodes: Node[];
}

export default function App() {
  const [title, setTitle] = useState("CloudStatus");
  const [nodes, setNodes] = useState<Node[]>([]);
  const [error, setError] = useState<string | null>(null);
  const location = useLocation();

  useEffect(() => {
    fetch("/api/nodes")
      .then((response) => response.json())
      .then((data: ApiResponse) => {
        setTitle(data.title);
        setNodes(data.nodes);
      })
      .catch((err) => {
        setError("Failed to load navigation data");
        console.error("Error fetching nodes:", err);
      });
  }, []);

  return (
    <Layout title={title}>
      <Navbar bg="dark" variant="dark" expand="lg">
        <Container>
          <Navbar.Brand as={Link} to="/">{title}</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="me-auto">
              <Nav.Link as={Link} to="/" active={location.pathname === "/"}>概览</Nav.Link>
              {nodes.map((node) => (
                <>
                  <Nav.Link 
                    key={`${node.id}`} 
                    as={Link} 
                    to={`/${node.id}`}
                    active={location.pathname === `/${node.id}`}
                  >{node.label}</Nav.Link>
                </>
              ))}
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
      {error && (
        <div className="alert alert-danger" role="alert">
          {error}
        </div>
      )}
      <Outlet />
    </Layout>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
