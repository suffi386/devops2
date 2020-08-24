import { Injectable } from '@angular/core';

import { StorageService } from '../storage.service';
import { GrpcInterceptor } from './grpc-interceptor';

const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer ';
const accessTokenStorageField = 'access_token';

@Injectable({ providedIn: 'root' })
export class GrpcAuthInterceptor implements GrpcInterceptor {
    constructor(private readonly authStorage: StorageService) { }

    public async intercept(
        request: any,
        invoker: any,
    ): Promise<any> {
        console.log(request);
        const reqMsg = request.getRequestMessage();
        reqMsg.setMessage('[Intercept request]' + reqMsg.getMessage());

        // After the RPC returns successfully, update the response.
        return invoker(request).then((response: any) => {
            // You can also do something with response metadata here.
            console.log(response.getMetadata());

            // Update the response message.
            const responseMsg = response.getResponseMessage();

            const accessToken = this.authStorage.getItem(accessTokenStorageField);
            if (accessToken) {
                responseMsg.setMetadata({ [authorizationKey]: bearerPrefix + accessToken });
            }

            responseMsg.setMessage('[Intercept response]' + responseMsg.getMessage());

            return response;
        });
    }
}
