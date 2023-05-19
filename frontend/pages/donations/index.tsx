import Dropdown from "@components/Dropdown";
import FilterDropdown from "@components/FilterDropdown";
import Layout from "@components/Layout";
import { useEffect, useRef, useState } from "react";
import allTags from "lib/tag-types";
import { BiSearch } from "react-icons/bi";
import Donation from "lib/types/donation";
import { GetStaticProps } from "next";
import RawDonation from "lib/types/rawDonation";
import axios from "axios";

const sortByOptions = [
    {
        name: "Sort by",
        func: null,
        id: 1
    },
    {
        name: "Date Added (asc.)",
        func: (a: Donation, b: Donation) => (a.creation_timestamp.getTime() - b.creation_timestamp.getTime()),
        id: 2
    },
    {
        name: "Date Added (desc.)",
        func: (a: Donation, b: Donation) => -(a.creation_timestamp.getTime() - b.creation_timestamp.getTime()),
        id: 3
    },
    {
        name: "Alphabetical (asc.)",
        func: (a: Donation, b: Donation) => a.title.localeCompare(b.title),
        id: 4
    },
    {
        name: "Alphabetical (desc.)",
        func: (a: Donation, b: Donation) => -a.title.localeCompare(b.title),
        id: 5
    },
]

const filterByDateOptions = {
    1: {
        name: "< 1 day ago",
        timeDelta: 1000 * 60 * 60 * 24, // 24h/day * 60min/h * 60s/min * 1000ms/s
        id: 1
    },
    2: {
        name: "< 1 week ago",
        timeDelta: 1000 * 60 * 60 * 24 * 7, // 7days/week * 24h/day * 60min/h * 60s/min * 1000ms/s
        id: 2
    },
    3: {
        name: "< 2 weeks ago",
        timeDelta: 1000 * 60 * 60 * 24 * 7 * 2, // 2 weeks * 7days/week * 24h/day * 60min/h * 60s/min * 1000ms/s
        id: 3
    },
    4: {
        name: "< 1 month ago",
        timeDelta: 1000 * 60 * 60 * 24 * 7 * 4, // 4weeks/month * 7days/week * 24h/day * 60min/h * 60s/min * 1000ms/s
        id: 4
    },
}

const tagsOptions = allTags.reduce((a, v) => ({ ...a, [v.id]: v }), {})

export const getStaticProps: GetStaticProps = async (context) => {
    const rawDonations: RawDonation[] = (await axios.get("http://localhost:8080/donations/list")).data

    const props = { rawDonations }
    return { props, revalidate: 1 }
}

export default function DonationsIndex({ rawDonations }: { rawDonations: RawDonation[] }) {
    const originalDonations: Donation[] = rawDonations.map(rawDonation => ({
        ...rawDonation,
        creation_timestamp: new Date(rawDonation.creation_timestamp)
    }))

    const searchBoxRef = useRef<HTMLInputElement>()
    const [data, setData] = useState<Donation[]>(originalDonations)
    const [sortBy, setSortBy] = useState(sortByOptions[0])
    const [filterByDate, setFilterByDate] = useState<number[]>([]) // These are arrays of IDs, not objects
    const [filterByTags, setFilterByTags] = useState<number[]>([])

    const applySortAndFilter = () => {
        // First filter by query
        const searchQuery = searchBoxRef.current.value;
        const filteredByQuery = originalDonations.filter(donation => donation.title.toLocaleLowerCase().includes(searchQuery.toLocaleLowerCase()));

        // Next filter by tags and date
        const tagsToFilterBy = filterByTags.map(id => tagsOptions[id].name);
        // @ts-ignore This will always be a number
        let largestTimeDelta = Math.max(filterByDate.map(id => filterByDateOptions[id].timeDelta));
        largestTimeDelta = largestTimeDelta === 0 ? 10e10 : largestTimeDelta;

        const filteredByTagsAndTime = filteredByQuery.filter(donation => (
            (tagsToFilterBy.length !== 0 ? tagsToFilterBy.some(tag => donation.tags.includes(tag)) : true) &&
            Date.now() - donation.creation_timestamp.getTime() <= largestTimeDelta
        ))

        // Finally, sort by query. If they don't want to sort by anything, just return the filtered version.
        const finalData = sortBy.func !== null ? [...filteredByTagsAndTime].sort(sortBy.func) : filteredByTagsAndTime;

        // Finally set the state to the filtered data
        setData(finalData)
    }

    // Make sure to re-filter and sort every time the criteria changes
    useEffect(() => {
        applySortAndFilter();
    }, [sortBy, filterByDate, filterByTags])

    return (
        <Layout name="Donations">
            <div className="flex flex-col gap-4 p-10 flex-grow min-w-screen">
                <div className="flex flex-col gap-2">
                    <h1 className="text-4xl text-white font-bold text-center">Donation Search</h1>
                    <p className="text-md text-gray-200 text-center">Search to find donation offers.</p>
                </div>

                <div className="flex flex-col self-center gap-2 w-full md:w-2/3 lg:w-1/2">
                    <div className="flex flex-col gap-1 flex-grow">
                        <label htmlFor="search">
                            <p className="text-sm text-gray-500 font-semibold">
                                Query
                            </p>
                        </label>

                        <div className="flex gap-2">
                            <input
                                className="bg-slate-700 appearance-none outline-none px-4 py-1 rounded focus:ring focus:bg-gray-600 duration-75 ring-0 w-full text-white placeholder:text-gray-400"
                                type="search"
                                placeholder="Ex. clothing"
                                id="search"
                                ref={searchBoxRef}
                            />
                            <button className="py-2 px-4 bg-blue-500 hover:bg-blue-600 active:bg-blue-700 duration-75 rounded-lg font-medium text-white" onClick={() => applySortAndFilter()}>
                                <BiSearch />
                            </button>
                        </div>

                        <div className="flex flex-wrap gap-2 self-center items-center justify-center">
                            <Dropdown title="Sort by" selectedItem={sortBy} setSelectedItem={setSortBy} options={sortByOptions} />
                            <FilterDropdown title="Filter by date" selectedItems={filterByDate} setSelectedItems={setFilterByDate} options={filterByDateOptions} />
                            <FilterDropdown title="Filter by tags" selectedItems={filterByTags} setSelectedItems={setFilterByTags} options={tagsOptions} />
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    )
}