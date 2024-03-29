/**
 * Data schema for a Donation.
 */
export default interface Donation {
    id: string,
    title: string,
    description: string, // markdown
    location: string,
    img: string, // direct src to firebase image
    creation_timestamp: Date,
    tags: string[] | null,
    owner_id: string,
    reports: string[]
}