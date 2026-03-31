

import * as Models from "./models";
export * from "./models";
import { TestingAPIError, ReasonTransport, ReasonEncoding, ReasonUnexpected } from "./models";

export const TestingAPIVersion = "1.0.0"

export class TestingAPI {
  private baseURL: string;
  private headers: Record<string, string>;
  private fetch: typeof fetch;

  constructor(baseURL: string, customFetch?: typeof fetch) {
    this.baseURL = baseURL;
    this.fetch = customFetch || fetch;
    this.headers = {
      'Accept': 'application/json',
      'User-Agent': 'TestingAPI-TypeScriptSDK/1.0.0'
    };
  }

  private addHeaders(request: RequestInit): RequestInit {
    request.headers = {...this.headers, ...request.headers };
    return request;
  }

  
  
  // Throws TestingAPIError, or a network error
  async CreateUser(params: Models.CreateUserReq): Promise<CreateUserResult> {
    var result = {} as CreateUserResult;

    var path = "/users/new";
    

    const url = new URL(path, this.baseURL);
    

    var requestInit: RequestInit = {
      method: "POST",
    };
    
    
    
    var authAdminToken = paramToString(params.AdminTokenAuth, "auth-header: X-App-Admin-Token", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-Admin-Token": authAdminToken};
    
    
    
    var authAPIKey = paramToString(params.APIKeyAuth, "auth-header: X-App-API-Key", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-API-Key": authAPIKey};
    
    
    
    
    requestInit.body = JSON.stringify(params.Body);
    requestInit.headers = { ...requestInit.headers, "Content-Type": "application/json"};
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 201:
        result.Response201 = await Models.ParseCreateUser201(response)
        break;
      
    
      
      case 400:
        result.Response400 = await Models.ParseCreateUser400(response)
        break;
      
    
      
      // CreateUser413 is a status-code only response
      // Payload Too Large - the request body exceeds the maximum allowed size
      
    
      
      case 500:
        result.Response500 = await Models.ParseCreateUser500(response)
        break;
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async GetUser(params: Models.GetUserReq): Promise<GetUserResult> {
    var result = {} as GetUserResult;

    var path = "/users/{userId}";
    
    var pathParamUserId = paramToString(params.UserId, "path parameter: userId", "string", true);
    path = path.replace("{userId}", encodeURIComponent(pathParamUserId));
    

    const url = new URL(path, this.baseURL);
    

    var requestInit: RequestInit = {
      method: "GET",
    };
    
    
    
    var authAPIKey = paramToString(params.APIKeyAuth, "auth-header: X-App-API-Key", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-API-Key": authAPIKey};
    
    
    
    var authSessionToken = paramToString(params.SessionTokenAuth, "auth-header: X-App-Session-Token", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-Session-Token": authSessionToken};
    
    
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 200:
        result.Response200 = await Models.ParseGetUser200(response)
        break;
      
    
      
      case 400:
        result.Response400 = await Models.ParseGetUser400(response)
        break;
      
    
      
      case 404:
        result.Response404 = await Models.ParseGetUser404(response)
        break;
      
    
      
      case 500:
        result.Response500 = await Models.ParseGetUser500(response)
        break;
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async ListUsers(params: Models.ListUsersReq): Promise<ListUsersResult> {
    var result = {} as ListUsersResult;

    var path = "/users";
    

    const url = new URL(path, this.baseURL);
    
    var queryParamPageNumber = paramToString(params.PageNumber, "query parameter: page", "integer", false);
    if (queryParamPageNumber != "") {
      url.searchParams.append("page", queryParamPageNumber);
    }
    
    var queryParamPageSize = paramToString(params.PageSize, "query parameter: pageSize", "integer", false);
    if (queryParamPageSize != "") {
      url.searchParams.append("pageSize", queryParamPageSize);
    }
    

    var requestInit: RequestInit = {
      method: "GET",
    };
    
    
    
    var authAdminToken = paramToString(params.AdminTokenAuth, "auth-header: X-App-Admin-Token", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-Admin-Token": authAdminToken};
    
    
    
    var authAPIKey = paramToString(params.APIKeyAuth, "auth-header: X-App-API-Key", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-API-Key": authAPIKey};
    
    
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 200:
        result.Response200 = await Models.ParseListUsers200(response)
        break;
      
    
      
      case 400:
        result.Response400 = await Models.ParseListUsers400(response)
        break;
      
    
      
      case 500:
        result.Response500 = await Models.ParseListUsers500(response)
        break;
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async LogoutUser(params: Models.LogoutUserReq): Promise<LogoutUserResult> {
    var result = {} as LogoutUserResult;

    var path = "/users/logout";
    

    const url = new URL(path, this.baseURL);
    

    var requestInit: RequestInit = {
      method: "GET",
    };
    
    
    
    var authAPIKey = paramToString(params.APIKeyAuth, "auth-header: X-App-API-Key", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-API-Key": authAPIKey};
    
    
    
    var numAuthParamsSet = 0;
    
    
    try {
      var authRefreshToken = paramToString(params.RefreshTokenAuth, "auth-header: X-App-Refresh-Token", "string", true);
      requestInit.headers = {...requestInit.headers, "X-App-Refresh-Token": authRefreshToken};
      numAuthParamsSet++;
    } catch (e) {
      // do nothing if it's a TestingAPIError, since we mandated the setting of the param to check if it's given or not even though it's a set-one-of auth param.
      if (!(e instanceof TestingAPIError)) {
        throw e;
      }
    }
    
    
    
    try {
      var authSessionToken = paramToString(params.SessionTokenAuth, "auth-header: X-App-Session-Token", "string", true);
      requestInit.headers = {...requestInit.headers, "X-App-Session-Token": authSessionToken};
      numAuthParamsSet++;
    } catch (e) {
      // do nothing if it's a TestingAPIError, since we mandated the setting of the param to check if it's given or not even though it's a set-one-of auth param.
      if (!(e instanceof TestingAPIError)) {
        throw e;
      }
    }
    
    

    if (numAuthParamsSet != 1) {
      throw new TestingAPIError(ReasonEncoding, `Exactly one auth param must be set, but ${numAuthParamsSet} were set instead`);
    }
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 200:
        result.Response200 = await Models.ParseLogoutUser200(response)
        break;
      
    
      
      case 400:
        result.Response400 = await Models.ParseLogoutUser400(response)
        break;
      
    
      
      case 500:
        result.Response500 = await Models.ParseLogoutUser500(response)
        break;
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async WhoAmI(params: Models.WhoAmIReq, body: BodyInit): Promise<WhoAmIResult> {
    var result = {} as WhoAmIResult;

    var path = "/users/whoami";
    

    const url = new URL(path, this.baseURL);
    

    var requestInit: RequestInit = {
      method: "POST",
    };
    
    
    
    var authAPIKey = paramToString(params.APIKeyAuth, "auth-header: X-App-API-Key", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-API-Key": authAPIKey};
    
    
    
    var authSessionToken = paramToString(params.SessionTokenAuth, "auth-header: X-App-Session-Token", "string", true);
    requestInit.headers = {...requestInit.headers, "X-App-Session-Token": authSessionToken};
    
    
    
    
    requestInit.body = body;
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 200:
        result.Response200 = await Models.ParseWhoAmI200(response)
        break;
      
    
      
      // WhoAmI400 is a status-code only response
      // Invalid Request
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async HealthCheck(params: Models.HealthCheckReq): Promise<HealthCheckResult> {
    var result = {} as HealthCheckResult;

    var path = "/health";
    

    const url = new URL(path, this.baseURL);
    

    var requestInit: RequestInit = {
      method: "GET",
    };
    
    
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await this.fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      
      case 200:
        result.Response200 = await Models.ParseHealthCheck200(response)
        break;
      
    
      default:
        result.UnknownResponse = response;
        break;
    }
    return result;
  }
  
}


export type CreateUserResult = {
  StatusCode: number;
  
  
  Response201: Models.CreateUser201;
  
  
  
  Response400: Models.CreateUser400;
  
  
  
  
  
  Response500: Models.CreateUser500;
  
  
  UnknownResponse: Response;
};

export type GetUserResult = {
  StatusCode: number;
  
  
  Response200: Models.GetUser200;
  
  
  
  Response400: Models.GetUser400;
  
  
  
  Response404: Models.GetUser404;
  
  
  
  Response500: Models.GetUser500;
  
  
  UnknownResponse: Response;
};

export type ListUsersResult = {
  StatusCode: number;
  
  
  Response200: Models.ListUsers200;
  
  
  
  Response400: Models.ListUsers400;
  
  
  
  Response500: Models.ListUsers500;
  
  
  UnknownResponse: Response;
};

export type LogoutUserResult = {
  StatusCode: number;
  
  
  Response200: Models.LogoutUser200;
  
  
  
  Response400: Models.LogoutUser400;
  
  
  
  Response500: Models.LogoutUser500;
  
  
  UnknownResponse: Response;
};

export type WhoAmIResult = {
  StatusCode: number;
  
  
  Response200: Models.WhoAmI200;
  
  
  
  
  UnknownResponse: Response;
};

export type HealthCheckResult = {
  StatusCode: number;
  
  
  Response200: Models.HealthCheck200;
  
  
  UnknownResponse: Response;
};


function paramToString(param: any, paramDescription: string, expectedType: string, required: boolean): string {
  if (param === undefined || param === null) {
    if (required) {
      throw new TestingAPIError(ReasonEncoding, `${paramDescription} is required but was not provided`);
    } else {
      return "";
    }
  }

  if (typeof param === "string") {
    if (required && param.trim() === "") {
      throw new TestingAPIError(ReasonEncoding, `${paramDescription} is required but was not provided`);
    }
    return param;
  } else if (typeof param === "number" || typeof param === "boolean") {
    if (expectedType == "int") {
      var i = parseInt(param.toString());
      if (isNaN(i)) {
        throw new TestingAPIError(ReasonEncoding, `${paramDescription} should be an integer but got ${param}`);
      }
      return i.toString();
    } else if (expectedType == "double") {
      var d = parseFloat(param.toString());
      if (isNaN(d)) {
        throw new TestingAPIError(ReasonEncoding, `${paramDescription} should be a number but got ${param}`);
      }
      return d.toString();
    }
    return param.toString();
  } else {
    throw new TestingAPIError(ReasonEncoding, `${paramDescription} should be of type ${expectedType} but got ${typeof param}`);
  }
}
