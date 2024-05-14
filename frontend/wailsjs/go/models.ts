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
	
	export class Oath {
	    Id: number;
	    CreationDate: string;
	    // Go type: time
	    StagedforFiling: any;
	    SubmitterName: string;
	
	    static createFrom(source: any = {}) {
	        return new Oath(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.CreationDate = source["CreationDate"];
	        this.StagedforFiling = this.convertValues(source["StagedforFiling"], null);
	        this.SubmitterName = source["SubmitterName"];
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

}

