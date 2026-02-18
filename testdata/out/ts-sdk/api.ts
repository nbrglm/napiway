

import { TestingAPIError } from "./models";
import * as Models from "./models";

export class TestingAPI {
  private baseURL: string;
  private headers: Record<string, string>;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    this.headers = {
      'Accept': 'application/json',
      'User-Agent': 'TestingAPI-TypeScriptSDK/1.0.0'
    };
  }

  private addHeaders(request: RequestInit): RequestInit {
    request.headers = { ...this.headers, ...request.headers };
    return request;
  }

  
  
  // Throws TestingAPIError, or a network error
  async CreateUser(params: Models.CreateUserRequest): Promise<CreateUserResult> {
    var result = {} as CreateUserResult;
    Models.ValidateCreateUserRequest(params);

    var path = "/users/new";
    
    const url = new URL(this.baseURL + path);
    

    var requestInit: RequestInit = {
      method: "POST",
    }
    
    
    requestInit.body = JSON.stringify(params.Body);
    requestInit.headers = { ...requestInit.headers, 'Content-Type': 'application/json' };
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      case 201:
        result.Response201 = await Models.NewCreateUser201Response(response);
        break;
    
      case 400:
        result.Response400 = await Models.NewCreateUser400Response(response);
        break;
    
      case 500:
        result.Response500 = await Models.NewCreateUser500Response(response);
        break;
    
      default:
        result.UnknownResponse = {
          StatusCode: response.status,
          Resp: response,
        };
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async GetUser(params: Models.GetUserRequest): Promise<GetUserResult> {
    var result = {} as GetUserResult;
    Models.ValidateGetUserRequest(params);

    var path = "/users/{userId}";
    
    var pathParamUserId = paramToString(params.UserId, "path parameter: UserId", "string", true);
    path = path.replace("{userId}", encodeURIComponent(pathParamUserId));
    
    const url = new URL(this.baseURL + path);
    

    var requestInit: RequestInit = {
      method: "GET",
    }
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      case 200:
        result.Response200 = await Models.NewGetUser200Response(response);
        break;
    
      case 400:
        result.Response400 = await Models.NewGetUser400Response(response);
        break;
    
      case 404:
        result.Response404 = await Models.NewGetUser404Response(response);
        break;
    
      case 500:
        result.Response500 = await Models.NewGetUser500Response(response);
        break;
    
      default:
        result.UnknownResponse = {
          StatusCode: response.status,
          Resp: response,
        };
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async HealthCheck(params: Models.HealthCheckRequest): Promise<HealthCheckResult> {
    var result = {} as HealthCheckResult;
    Models.ValidateHealthCheckRequest(params);

    var path = "/health";
    
    const url = new URL(this.baseURL + path);
    

    var requestInit: RequestInit = {
      method: "GET",
    }
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      case 200:
        result.Response200 = await Models.NewHealthCheck200Response(response);
        break;
    
      default:
        result.UnknownResponse = {
          StatusCode: response.status,
          Resp: response,
        };
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async ListUsers(params: Models.ListUsersRequest): Promise<ListUsersResult> {
    var result = {} as ListUsersResult;
    Models.ValidateListUsersRequest(params);

    var path = "/users";
    
    const url = new URL(this.baseURL + path);
    
    var queryParamPageNumber = paramToString(params.PageNumber, "query parameter: PageNumber", "number", false);
    if (queryParamPageNumber !== "") {
      url.searchParams.append("page", queryParamPageNumber);
    }
    
    var queryParamPageSize = paramToString(params.PageSize, "query parameter: PageSize", "number", false);
    if (queryParamPageSize !== "") {
      url.searchParams.append("pageSize", queryParamPageSize);
    }
    

    var requestInit: RequestInit = {
      method: "GET",
    }
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      case 200:
        result.Response200 = await Models.NewListUsers200Response(response);
        break;
    
      case 400:
        result.Response400 = await Models.NewListUsers400Response(response);
        break;
    
      case 500:
        result.Response500 = await Models.NewListUsers500Response(response);
        break;
    
      default:
        result.UnknownResponse = {
          StatusCode: response.status,
          Resp: response,
        };
    }
    return result;
  }
  
  
  // Throws TestingAPIError, or a network error
  async LogoutUser(params: Models.LogoutUserRequest): Promise<LogoutUserResult> {
    var result = {} as LogoutUserResult;
    Models.ValidateLogoutUserRequest(params);

    var path = "/users/logout";
    
    const url = new URL(this.baseURL + path);
    

    var requestInit: RequestInit = {
      method: "GET",
    }
    
    
    const request = new Request(url, this.addHeaders(requestInit));
    const response = await fetch(request);
    result.StatusCode = response.status;
    switch (response.status) {
    
      case 200:
        result.Response200 = await Models.NewLogoutUser200Response(response);
        break;
    
      case 400:
        result.Response400 = await Models.NewLogoutUser400Response(response);
        break;
    
      case 500:
        result.Response500 = await Models.NewLogoutUser500Response(response);
        break;
    
      default:
        result.UnknownResponse = {
          StatusCode: response.status,
          Resp: response,
        };
    }
    return result;
  }
  
}


type CreateUserResult = {
  StatusCode: number;
  
  Response201: Models.CreateUser201Response;
  
  Response400: Models.CreateUser400Response;
  
  Response500: Models.CreateUser500Response;
  
  UnknownResponse: UnknownStatusResponse;
}

type GetUserResult = {
  StatusCode: number;
  
  Response200: Models.GetUser200Response;
  
  Response400: Models.GetUser400Response;
  
  Response404: Models.GetUser404Response;
  
  Response500: Models.GetUser500Response;
  
  UnknownResponse: UnknownStatusResponse;
}

type HealthCheckResult = {
  StatusCode: number;
  
  Response200: Models.HealthCheck200Response;
  
  UnknownResponse: UnknownStatusResponse;
}

type ListUsersResult = {
  StatusCode: number;
  
  Response200: Models.ListUsers200Response;
  
  Response400: Models.ListUsers400Response;
  
  Response500: Models.ListUsers500Response;
  
  UnknownResponse: UnknownStatusResponse;
}

type LogoutUserResult = {
  StatusCode: number;
  
  Response200: Models.LogoutUser200Response;
  
  Response400: Models.LogoutUser400Response;
  
  Response500: Models.LogoutUser500Response;
  
  UnknownResponse: UnknownStatusResponse;
}


type UnknownStatusResponse = {
  StatusCode: number;
  Resp: Response;
}

function paramToString(param: any, paramDescription: string, expectedType: string, required: boolean): string {
  if (param === undefined || param === null) {
    if (required) {
      throw new TestingAPIError(Models.TestingAPIErrorReasonInvalidRequest, `${paramDescription} is required but was not provided`);
    } else {
      return "";
    }
  }

  if (typeof param === "string") {
    return param;
  } else if (typeof param === "number" || typeof param === "boolean") {
    return param.toString();
  } else {
    throw new TestingAPIError(Models.TestingAPIErrorReasonInvalidRequest, `${paramDescription} should be of type ${expectedType} but got ${typeof param}`);
  }
}

