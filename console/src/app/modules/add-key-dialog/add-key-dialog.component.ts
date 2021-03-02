import { Component, Inject } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

export enum AddKeyDialogType {
    MACHINE = "MACHINE",
    AUTHNKEY = "AUTHNKEY",
}

@Component({
    selector: 'app-add-key-dialog',
    templateUrl: './add-key-dialog.component.html',
    styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
    public startDate: Date = new Date();
    types: MachineKeyType[] | AuthNKeyType[] = [];
    public type!: MachineKeyType | AuthNKeyType;
    public dateControl: FormControl = new FormControl('', []);

    constructor(
        public dialogRef: MatDialogRef<AddKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        if (data.type = AddKeyDialogType.MACHINE) {
            this.types = [MachineKeyType.MACHINEKEY_JSON];
            this.type = MachineKeyType.MACHINEKEY_JSON;
        } else if (data.type = AddKeyDialogType.AUTHNKEY) {
            this.types = [AuthNKeyType.AUTHNKEY_JSON];
            this.type = AuthNKeyType.AUTHNKEY_JSON;
        }
        const today = new Date();
        this.startDate.setDate(today.getDate() + 1);
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ type: this.type, date: this.dateControl.value });
    }
}
