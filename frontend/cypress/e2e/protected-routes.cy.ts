describe('Protected Routes', () => {
  it('should not show up in the navigation bar', () => {
    // Start from the index page
    cy.visit('http://localhost:3000');
    cy.get("header").contains("Home"); // base case
    cy.get("header").contains("Donate").should("not.exist");
    cy.get("header").contains("Profile").should("not.exist");
  });

  it('should redirect me back if I am not logged in', () => {
    cy.on('window:alert', (txt) => {
      // Mocha assertions
      expect(txt).to.contains('You need to be signed in to access this page. Redirecting...');
    })

    cy.visit('http://localhost:3000/profile');
    cy.wait(1000); // Wait for alert to show up
    cy.visit('http://localhost:3000/donations/create');
  });

});