import { Fragment } from 'react'
import { Menu, Transition } from '@headlessui/react'
import { BsChevronDown } from "react-icons/bs"

/**
 * Converts a list of class names into one className string
 * @param classes List of classes
 * @returns className string
 */
function classNames(...classes) {
    return classes.filter(Boolean).join(' ')
}

/**
 * A customizable Dropdown that updates state of a parent component.
 * @param props All the props of the component. 
 */
export default function Dropdown({ title, selectedItem, setSelectedItem, options, openOverlap }: PropTypes) {
    // Default to True if not provided
    openOverlap = typeof openOverlap === "boolean" ? openOverlap : false;

    return (
        <Menu as="div" className="relative inline-block text-left w-36">
            <div>
                <Menu.Button className="inline-flex w-full justify-center gap-x-1.5 rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                    {selectedItem ? selectedItem.name : title}
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
                <Menu.Items className={classNames("right-0 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-[999999999999999999]", openOverlap ? "absolute" : "")}>
                    <div className="py-1">
                        {options.map(item => (
                            <Menu.Item key={item.name}>
                                {({ active }) => {
                                    return (
                                        <button
                                            onClick={() => setSelectedItem(item)}
                                            className={classNames(
                                                active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                                                'block px-4 py-2 text-sm w-full text-left'
                                            )}
                                        >
                                            {item.name}
                                        </button>
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
    selectedItem: {[key: string]: any},
    setSelectedItem: Function,
    options: any[],
    openOverlap?: Boolean
}