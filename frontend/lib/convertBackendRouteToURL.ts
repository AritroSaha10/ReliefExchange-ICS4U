export default function convertBackendRouteToURL(route: string) {
    return process.env.NEXT_PUBLIC_BACKEND_BASE_URL + route;
}