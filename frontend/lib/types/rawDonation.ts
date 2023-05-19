export default interface RawDonation {
    id: string,
    title: string,
    description: string, // markdown
    location: string,
    img: string, // direct src to firebase image
    creation_timestamp: string,
    tags: string[],
    owner_id: string,
    reports: string[]
}