/**
 * @file DonationCard component definition.
 * This file contains the implementation of the DonationCard component, which is used to 
 * display brief information about a donation. It includes the title, date, image (optional), 
 * subtitle, tags, donation link, and additional functionalities for administrators. 
 * @author Aritro Saha
 * @cite Installation - Tailwind CSS, https://tailwindcss.com/docs/installation. 
 */

import Image from 'next/image'
import Link from 'next/link'

import removeMd from "remove-markdown"
import { FiFlag } from 'react-icons/fi';

import DonationTag from '@lib/types/tag';

/**
 * A Donation card that shows brief information about a donation. 
 * @returns Donation component with data provided.
 */
export default function DonationCard({ title, date, image, subtitle, tags, href, isAdmin, reportCount }: { title: string, date: Date, image?: string, subtitle: string, tags: DonationTag[], href: string, isAdmin: boolean, reportCount: number }) {
    return (
        <div className="flex flex-col md:flex-row md:justify-start items-center gap-4 rounded-xl p-2 lg:p-6 duration-300 bg-slate-700">
            {image && <Image src={image} className="rounded-lg z-10 bg-blue-200 object-center object-cover aspect-square" alt={title} width={200} height={150} />}
            <div className="flex flex-col items-center md:items-start gap-2 lg:gap-4 w-full">
                <div className="flex flex-col-reverse md:flex-row justify-center md:justify-between w-full items-center gap-1 lg:gap-2">
                    <div>
                        <h1 className="text-xl text-white text-center md:text-left font-semibold">{title}</h1>
                        <div className='flex flex-wrap items-center gap-1 text-sm text-gray-200'>
                            <span className='text-center md:text-left'>
                                {date.toLocaleDateString("en-CA", {
                                    day: "numeric",
                                    month: "short",
                                    year: "numeric"
                                })}
                            </span>

                            {tags.length !== 0 && <span> | </span>}

                            <div className='flex flex-wrap gap-1'>
                                {tags.map(tag => (
                                    <span className={`text-xs text-white border-0 px-2 py-1 rounded-lg text-center ${tag.color}`} key={tag.name}>{tag.name.toLowerCase()}</span>
                                ))}
                            </div>

                            {isAdmin && (
                                <>
                                    <span> | </span>
                                    <div className='flex items-center text-red-500'>
                                        <FiFlag className='mr-1'></FiFlag>
                                        {reportCount}
                                    </div>
                                </>
                            )}
                        </div>
                    </div>
                </div>

                <p className="break-words text-center md:text-left text-white">
                    {removeMd(subtitle).slice(0, subtitle.length > 120 ? 120 : subtitle.length)}{subtitle.length > 120 && "..."}
                </p>

                <Link href={href} passHref>
                    <button className="text-sm bg-blue-500 hover:bg-blue-700 text-white font-bold w-24 h-8 rounded duration-300 mt-2">View More</button>
                </Link>
            </div>
        </div>
    );
}