import { NextFetchEvent, NextResponse, URLPattern } from 'next/server';
import type { NextRequest } from 'next/server';

// Gets the path segments from an entire path
const PATTERNS = [
    [
        new URLPattern({ pathname: '/__/auth/:path' }),
        ({ pathname }) => pathname.groups,
    ],
]

// Runs path extractors on url 
const params = (url) => {
    const input = url.split('?')[0]
    let result: {[key: string]: any} = {}

    for (const [pattern, handler] of PATTERNS) {
        // @ts-ignore
        const patternResult = pattern.exec(input)
        if (patternResult !== null && 'pathname' in patternResult) {
            // @ts-ignore
            result = handler(patternResult)
            break
        }
    }
    return result
}

// Serves the middleware properly
export async function middleware(request: NextRequest, event: NextFetchEvent) {
    const { path } = params(request.url);
    const baseUrl = request.url.substring(0, request.url.indexOf(new URL(request.url).pathname));

    const data = await (await fetch(new URL(`auth_helpers/${path}`, baseUrl))).text();
    const response = new NextResponse(data);
    response.headers.append("Content-Type", (path.includes(".js") ? "text/javascript;" : "text/html;") + " charset=utf-8")

    return response;
}

// See "Matching Paths" below to learn more
export const config = {
    matcher: '/__/auth/:path*',
};