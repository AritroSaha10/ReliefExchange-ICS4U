/**
 * Increase the resolution of a Google profile picture link. 
 * This is necessary since it's always 96px x 96px with no way to
 * request a higher resolution picture.
 * @param url The original image src link
 * @param newRes The new resolution in pixels (2-dimensional res not allowed)
 * @returns An altered src link with higher resolution
 */
const increasePFPResolution = (url: string, newRes: Number) => (
    url.replace("s96-c", `s${newRes.toFixed(0)}-c`)
)

export default increasePFPResolution;