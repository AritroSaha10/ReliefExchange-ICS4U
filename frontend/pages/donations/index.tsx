/**
 * @file File for the the donation index page, which is accessible at /donations.
 * More info about the component can be seen in its own dcumentation.
 * @author Aritro Saha
 */

import { useEffect, useRef, useState } from "react";
import { GetStaticProps } from "next";

import { onAuthStateChanged } from "firebase/auth";
import axios from "axios";
import { BiSearch } from "react-icons/bi";

import Dropdown from "@components/Dropdown";
import FilterDropdown from "@components/FilterDropdown";
import Layout from "@components/Layout";
import DonationCard from "@components/DonationCard";

import allTags from "@lib/tag-types";
import Donation from "@lib/types/donation";
import RawDonation from "@lib/types/rawDonation";
import convertBackendRouteToURL from "@lib/convertBackendRouteToURL";
import auth from "@lib/firebase/auth";

// All the options that the user can sort by, complete with functions
// for extensibility
const sortByOptions = [
    {
        name: "Sort by",
        func: (a: Donation, b: Donation) => -(a.creation_timestamp.getTime() - b.creation_timestamp.getTime()), // Function to use to sort, default to date desc.
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

// All the options that admins can sort by, complete with functions
// for extensibility
const adminSortByOptions = [
    {
        name: "Reports (asc.)",
        func: (a: Donation, b: Donation) => (a.reports ? a.reports.length : 0) - (b.reports ? b.reports.length : 0),
        id: 6
    },
    {
        name: "Reports (desc.)",
        func: (a: Donation, b: Donation) => (b.reports ? b.reports.length : 0) - (a.reports ? a.reports.length : 0),
        id: 7
    },
]

// All the date options the user can filter by
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

// All the tag options the user can filter by, retrieved by converting the tags array
// into a dict
const tagsOptions = allTags.reduce((a, v) => ({ ...a, [v.id]: v }), {})

/**
 * Part of Next.js, fetches all donation data.
 */
export const getStaticProps: GetStaticProps = async (context) => {
    // Request the backend for the donation list
    const rawDonations: RawDonation[] = (await axios.get(convertBackendRouteToURL("/donations/list"))).data

    const props = { rawDonations }
    return { props, revalidate: 1 } // Revalidate the data cache 1s after page load
}

/**
 * Component for the donations index page, which shows a list of all the donations posted
 * on the website. It also allows users to search, sort, and filter by certain attributes.
 */
export default function DonationsIndex({ rawDonations }: { rawDonations: RawDonation[] }) {
    // Convert the ISO string timestamps in the raw donations to Date objects
    const originalDonations: Donation[] = rawDonations.map(rawDonation => ({
        ...rawDonation,
        creation_timestamp: new Date(rawDonation.creation_timestamp)
    }))

    // Necessary state and ref hooks
    const searchBoxRef = useRef<HTMLInputElement>()
    const [data, setData] = useState<Donation[]>(originalDonations)
    const [sortBy, setSortBy] = useState(sortByOptions[0])
    const [filterByDate, setFilterByDate] = useState<number[]>([]) // These are arrays of IDs, not objects
    const [filterByTags, setFilterByTags] = useState<number[]>([])
    const [isAdmin, setIsAdmin] = useState<boolean>(false);

    /**
     * Refresh user-specific (whether they're admin) data on auth change
     */
    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, user => {
            // Only run if user is signed in
            if (user && Object.keys(user).length !== 0) {
                // Attempt to get the user's data from the backend
                axios.get(convertBackendRouteToURL(`/users/${user.uid}`)).then(async res => {
                    // Get admin attribute from data
                    setIsAdmin(res.data && !!res.data.admin);
                }).catch(err => {
                    // Silently log error
                    console.error(err);
                })
            }
        })

        return () => unsubscribe();
    }, []) // eslint-disable-line react-hooks/exhaustive-deps

    /**
     * Apply the sorting and filtering criteria on the original data from API.
     */
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
            (tagsToFilterBy.length !== 0 ? tagsToFilterBy.some(tag => donation.tags && donation.tags.includes(tag)) : true) &&
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
                                onKeyDown={e => {
                                    if (e.key === "Enter") {
                                        applySortAndFilter()
                                    }
                                }}
                            />
                            <button className="py-2 px-4 bg-blue-500 hover:bg-blue-600 active:bg-blue-700 duration-75 rounded-lg font-medium text-white" onClick={() => applySortAndFilter()}>
                                <BiSearch />
                            </button>
                        </div>

                        <div className="flex flex-wrap gap-2 self-center justify-center">
                            <Dropdown title="Sort by" selectedItem={sortBy} setSelectedItem={setSortBy} options={isAdmin ? sortByOptions.concat(adminSortByOptions) : sortByOptions} openOverlap={true} />
                            <FilterDropdown title="Filter by date" selectedItems={filterByDate} setSelectedItems={setFilterByDate} options={filterByDateOptions} />
                            <FilterDropdown title="Filter by tags" selectedItems={filterByTags} setSelectedItems={setFilterByTags} options={tagsOptions} />
                        </div>
                    </div>
                </div>

                <div className="flex flex-col self-center gap-4 lg:gap-6 w-full px-4 md:px-8 py-4 lg:px-12">
                    {data.map(donation => {
                        const tags = (
                            donation.tags ? donation.tags.map(tagName => allTags.find(tag => tag.name === tagName)) : []
                        ).filter(tag => tag !== undefined);

                        return (
                            <DonationCard
                                title={donation.title}
                                date={donation.creation_timestamp}
                                subtitle={donation.description}
                                image={donation.img}
                                tags={tags}
                                href={`/donations/${donation.id}`}
                                isAdmin={isAdmin}
                                reportCount={donation.reports ? donation.reports.length : 0}
                                key={donation.id}
                            />
                        )
                    })}

                    {data.length === 0 && (
                        <div className="flex flex-col items-center text-center">
                            <h2 className="text-4xl text-white font-semibold mb-2">No Results Found</h2>
                            <p className="text-xl text-gray-200 lg:w-3/4">
                                Sorry, but we couldn&apos;t find any donations matching your search query. 
                                Please adjust it accordingly or try again later. If you just created a donation 
                                but it isn&apos;t showing up, please refresh after a few seconds.
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </Layout>
    )
}