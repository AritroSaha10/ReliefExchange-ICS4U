import { render, screen } from "@testing-library/react";
import convertBackendRoundToURL from "../lib/convertBackendRouteToURL"

it("converts routes correctly", () => {
    const segments = "/this/is/a/route?a=woah"

    const result = convertBackendRoundToURL(segments)

    expect(result).toBe(process.env.NEXT_PUBLIC_BACKEND_BASE_URL + segments)
})
