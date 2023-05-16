import Layout from "@components/Layout";

export default function DonationsIndex() {
    return (
        <Layout name="Donations">
            <div className="flex flex-col gap-4 p-10 flex-grow min-w-screen">
                <div className="flex flex-col gap-2">
                    <h1 className="text-4xl text-white font-bold text-center">Grant Search</h1>
                    <p className="text-md text-gray-200 text-center">Search through a large choice of grants using our tools!</p>
                </div>

                <div className="flex flex-col self-center gap-2 w-full md:w-2/3 lg:w-1/2">
                    <div className="flex flex-col gap-1 flex-grow">
                        <label htmlFor="search">
                            <p className="text-sm text-gray-500 font-semibold">
                                Search
                            </p>
                        </label>

                        <input
                            className="bg-slate-700 appearance-none outline-none px-4 py-1 rounded focus:ring focus:bg-gray-600 duration-75 ring-0 w-full text-white placeholder:text-gray-400"
                            type="search"
                            placeholder="Ex. covid relief"
                            id="search"
                        />
                    </div>
                </div>
            </div>
        </Layout>
    )
}