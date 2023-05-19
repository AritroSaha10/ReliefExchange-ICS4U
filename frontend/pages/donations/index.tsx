import Dropdown from "@components/Dropdown";
import FilterDropdown from "@components/FilterDropdown";
import Layout from "@components/Layout";
import { useState } from "react";
import allTags from "lib/tag-types";

const sortByOptions = [
    {
        name: "Sort by",
        id: 1
    },
    {
        name: "Date Added (asc.)",
        id: 2
    },
    {
        name: "Date Added (desc.)",
        id: 3
    },
    {
        name: "Alphabetical (asc.)",
        id: 4
    },
    {
        name: "Alphabetical (desc.)",
        id: 5
    },
]

const filterByDateOptions = {
    1: {
        name: "< 1 day ago",
        id: 1
    },
    2: {
        name: "< 1 week ago",
        id: 2
    },
    3: {
        name: "< 2 weeks ago",
        id: 3
    },
    4: {
        name: "< 1 month ago",
        id: 4
    },
}

const tagsOptions = allTags.reduce((a, v) => ({ ...a, [v.id]: v }), {})

export default function DonationsIndex() {
    const [sortBy, setSortBy] = useState(sortByOptions[0])
    const [filterByDate, setFilterByDate] = useState([]) // These are arrays of IDs, not objects
    const [filterByTags, setFilterByTags] = useState([])

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
                                Search
                            </p>
                        </label>

                        <div className="flex gap-2">
                            <input
                                className="bg-slate-700 appearance-none outline-none px-4 py-1 rounded focus:ring focus:bg-gray-600 duration-75 ring-0 w-full text-white placeholder:text-gray-400"
                                type="search"
                                placeholder="Ex. clothing"
                                id="search"
                            />
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