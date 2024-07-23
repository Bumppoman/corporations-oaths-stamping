export namespace api {
	
	export class UserInfo {
	    Email: string;
	    Id: number;
	    IsHiddenInUI: boolean;
	    IsSiteAdmin: boolean;
	    LoginName: string;
	    PrincipalType: number;
	    Title: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Email = source["Email"];
	        this.Id = source["Id"];
	        this.IsHiddenInUI = source["IsHiddenInUI"];
	        this.IsSiteAdmin = source["IsSiteAdmin"];
	        this.LoginName = source["LoginName"];
	        this.PrincipalType = source["PrincipalType"];
	        this.Title = source["Title"];
	    }
	}

}

export namespace main {
	
	export class SignInResponse {
	    CanAccess: boolean;
	    CurrentVersion: string;
	    UserInfo?: api.UserInfo;
	
	    static createFrom(source: any = {}) {
	        return new SignInResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CanAccess = source["CanAccess"];
	        this.CurrentVersion = source["CurrentVersion"];
	        this.UserInfo = this.convertValues(source["UserInfo"], api.UserInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StampingItem {
	    Id: number;
	    Created: string;
	    Selected: boolean;
	    StampText: string;
	    Title: string;
	
	    static createFrom(source: any = {}) {
	        return new StampingItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Created = source["Created"];
	        this.Selected = source["Selected"];
	        this.StampText = source["StampText"];
	        this.Title = source["Title"];
	    }
	}

}

