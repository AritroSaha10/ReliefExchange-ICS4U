import { render, screen } from "@testing-library/react";
import '@testing-library/jest-dom'
import Home from "../pages/index"

it("renders homepage unchanged", () => {
    const { container } = render(<Home />)
    expect(container).toMatchSnapshot()
})

it("has a heading with the title", () => {
    render(<Home />)
    
    const heading = screen.getByRole("heading", {
        name: "ReliefExchange"
    })

    expect(heading).toBeInTheDocument()
})
