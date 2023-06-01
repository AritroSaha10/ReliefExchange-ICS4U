/**
 * Converts a list of class names into one className string
 * @param classes List of classes
 * @returns className string
 */
export default function classNames(...classes) {
    return classes.filter(Boolean).join(' ')
}