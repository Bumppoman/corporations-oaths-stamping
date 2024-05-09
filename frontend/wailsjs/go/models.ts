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
	    StagedforFiling: boolean;
	    SubmitterName: string;
	
	    static createFrom(source: any = {}) {
	        return new Oath(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.CreationDate = source["CreationDate"];
	        this.StagedforFiling = source["StagedforFiling"];
	        this.SubmitterName = source["SubmitterName"];
	    }
	}

}

