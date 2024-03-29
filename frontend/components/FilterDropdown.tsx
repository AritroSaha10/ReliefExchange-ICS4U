/**
 * @file FilterDropdown component definition.
 * 
 * This file contains the implementation of the FilterDropdown component, which is a customizable
 * dropdown menu that updates the state of a parent component, and is designed specifically for 
 * filtering data. It utilizes the '@headlessui/react' and 'react-icons/bs' modules for the menu 
 * and chevron icons respectively. The component accepts props for title, selectedItems, setSelectedItems, 
 * and options. The selectedItems prop represents the currently selected items, setSelectedItems is 
 * the function to update the selected items, and options is an object containing the available options 
 * for filtering. The types of each prop are also declared in PropTypes.
 * 
 * @author Aritro Saha
 */

import { Fragment } from 'react'

import { Menu, Transition } from '@headlessui/react'
import { BsChevronDown } from "react-icons/bs"

import classNames from "@lib/classNames"

/**
 * A customizable Dropdown that updates state of a parent component, specifically for
 * filtering data.
 * @param props All the props of the component. 
 */
export default function FilterDropdown({ title, selectedItems, setSelectedItems, options }: PropTypes) {
    return (
        <Menu as="div" className="relative inline-block text-left w-36">
            <div>
                <Menu.Button className="inline-flex w-full justify-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                    {title}
                    <BsChevronDown className="-mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
                </Menu.Button>
            </div>

            <Transition
                as={Fragment}
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
            >
                <Menu.Items className="absolute right-0 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-[999999999999999999]">
                    <div className="py-1">
                        {Object.values(options).map((item: any) => (
                            <Menu.Item key={item.id}>
                                {({ active }) => {
                                    return (
                                        <div className={classNames(
                                            active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                                            'flex gap-1 px-4 py-2 text-sm w-full text-left'
                                        )} key={item.id}>
                                            <input
                                                type="checkbox"
                                                onChange={(e) => {
                                                    if (e.target.checked) {
                                                        setSelectedItems([...selectedItems, item.id])
                                                    } else {
                                                        setSelectedItems(selectedItems.filter(a => a !== item.id))
                                                    }
                                                }}
                                                checked={selectedItems.includes(item.id)}
                                            />
                                            <p>{item.name}</p>
                                        </div>
                                    )
                                }}
                            </Menu.Item>
                        ))}
                    </div>
                </Menu.Items>
            </Transition>
        </Menu>
    )
}

interface PropTypes {
    title: string,
    selectedItems: any,
    setSelectedItems: Function,
    options: any,
}
