export default function convertBackendRouteToURL(route: string) {
    return process.env.RECAPTCHA_SECRET_KEY + route;
}