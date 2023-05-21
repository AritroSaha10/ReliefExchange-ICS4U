/**
 * Converts a route to the backend server to an actual URL
 * @param route backend route
 * @returns Full URL that requests can be sent to
 */
export default function convertBackendRouteToURL(route: string) {
    return process.env.NEXT_PUBLIC_BACKEND_BASE_URL + route;
}