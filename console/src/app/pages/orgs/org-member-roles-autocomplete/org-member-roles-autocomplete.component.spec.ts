import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OrgMemberRolesAutocompleteComponent } from './org-member-roles-autocomplete.component';

describe('OrgMemberRolesAutocompleteComponent', () => {
    let component: OrgMemberRolesAutocompleteComponent;
    let fixture: ComponentFixture<OrgMemberRolesAutocompleteComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [OrgMemberRolesAutocompleteComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgMemberRolesAutocompleteComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
