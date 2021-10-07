import { apiAuth } from "../../support/api/apiauth";
import { ensureProjectExists } from "../../support/api/projects";
import { login, User } from "../../support/login/users";

describe('permissions', () => {

    const testProjectName = 'e2eprojectpermission'
    const testRoleName = 'e2eroleundertest'
    const testGrantName = 'e2egrantundertest'

    ;[User.OrgOwner].forEach(user => {

        describe(`as user "${user}"`, () => {

            describe('add role', () => {

                before(()=> {
                    login(user)
                    apiAuth().then(api => {
                        ensureProjectExists(api, testProjectName).then(projectID => {
                            ensureRoleDoesntExist(api, projectID, testAppName).then(() => {
                                cy.visit(`${Cypress.env('consoleUrl')}/projects/${projectID}`)
                            })
                        })
                    })

    //                cy.consolelogin(`${user.toLowerCase()}_user_name@caos-demo.${Cypress.env('domain')}`, Cypress.env(`${user.toLowerCase()}_password`))
                    cy.visit(Cypress.env('consoleUrl') + '/projects')
                    // wait until table is loaded
                    cy.contains("tr", "cypress").contains("e2e")
                })
            })
        })
    })
})


describe('permissions', () => {

    before(()=> {
//        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })

    it('should show projects ', () => {
        cy.visit(Cypress.env('consoleUrl') + '/projects')
        cy.url().should('contain', '/projects')
    })

    it('should add a role', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")').click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
        });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/projects/' + projectID +'/roles/create'))
        cy.get('[formcontrolname^=key]').type("newdemorole")
        cy.get('[formcontrolname^=displayName]').type("newdemodisplayname")
        cy.get('[formcontrolname^=group]').type("newdemogroupname")
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the Role get processed
        cy.wait(5000)
    })

    it('should add a grant', () => {
        cy.visit(Cypress.env('consoleUrl') + '/org').then(() => {
            cy.url().should('contain', '/org');
        })
        cy.visit(Cypress.env('consoleUrl') + '/projects').then(() => {
            cy.url().should('contain', '/projects');
            cy.get('.card').should('contain.text', "newProjectToTest")
        })
        cy.get('.card').filter(':contains("newProjectToTest")').click()
        cy.get('.app-container').filter(':contains("newAppToTest")').should('be.visible').click()
        let projectID
        cy.url().then(url => {
            cy.log(url.split('/')[4])
            projectID = url.split('/')[4]
        });
        
        cy.then(() => cy.visit(Cypress.env('consoleUrl') + '/grant-create/project/' + projectID ))
        cy.get('input').type("demo")
        cy.get('[role^=listbox]').filter(`:contains("${Cypress.env("fullUserName")}")`).should('be.visible').click()
        cy.wait(5000)
        //cy.get('.button').contains('Continue').click()
        cy.get('button').filter(':contains("Continue")').click()
        cy.wait(5000)
        cy.get('tr').filter(':contains("demo")').find('label').click()
        cy.get('button').filter(':contains("Save")').should('be.visible').click()
        //let the grant get processed
        cy.wait(5000)
    })
})

