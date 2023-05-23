import { Component, OnInit } from '@angular/core';
import { PrivacyPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss'],
})
export class FooterComponent implements OnInit {
  public policy?: PrivacyPolicy.AsObject;
  constructor(public authService: GrpcAuthService) {}

  ngOnInit(): void {
    this.authService.getMyPrivacyPolicy().then((policyResp) => {
      if (policyResp.policy) {
        this.policy = policyResp.policy;
      }
    });
  }
}
