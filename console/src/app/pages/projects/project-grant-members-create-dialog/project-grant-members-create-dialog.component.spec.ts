import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectGrantMembersCreateDialogComponent } from './project-grant-members-create-dialog.component';

describe('ProjectGrantMembersCreateDialogComponent', () => {
    let component: ProjectGrantMembersCreateDialogComponent;
    let fixture: ComponentFixture<ProjectGrantMembersCreateDialogComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantMembersCreateDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantMembersCreateDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
