import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import { listFrameworks, hasFramework, getFramework } from '@netlify/framework-info';
import { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';

export type FrameworkDefinition = {
  id?: FrameworkName | string;
  title: string;
  imgSrcDark: string;
  imgSrcLight?: string;
  docsLink: string;
  external?: boolean;
};

export type Framework = FrameworkDefinition & {
  fragment: string;
};

@Component({
  standalone: true,
  selector: 'cnsl-quickstart',
  templateUrl: './quickstart.component.html',
  styleUrls: ['./quickstart.component.scss'],
  imports: [TranslateModule, RouterModule, CommonModule, MatButtonModule],
})
export class QuickstartComponent {
  public frameworks: FrameworkDefinition[] = frameworkDefinition.map((f) => {
    return {
      ...f,
      imgSrcDark: `assets${f.imgSrcDark}`,
      imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
    };
  });

  constructor() {
    // console.log(this.frameworks[0].title);
  }
}
