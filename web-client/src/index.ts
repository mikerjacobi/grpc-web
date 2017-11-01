import {grpc, Code, Metadata} from "grpc-web-client";
import {AccountService} from "../generated/app_pb_service";
import {Account, AuthenticateAccountResp} from "../generated/app_pb";

class AccountsAPI {
    host: string;
    currAccount: Account;

    constructor() { 
        this.host = "https://dev.jacobra.com:8443"
    }
    update(){
        var username = document.getElementById("current_username");
        var accountID = document.getElementById("current_account_id");
        username.innerHTML = this.currAccount.getUsername();
        accountID.innerHTML = this.currAccount.getAccountid();
    }
    create(){
        let aReq = new Account();
        aReq.setUsername((<HTMLInputElement>document.getElementById("create_username")).value);
        aReq.setPassword((<HTMLInputElement>document.getElementById("create_password")).value);

        grpc.invoke(AccountService.Create, {
            request: aReq,
            host: this.host,
            onMessage: (message: Account) => {
                this.currAccount = message;
			    this.update();	
            },
            onEnd: (code: Code, msg: string | undefined, trailers: Metadata) => {
                if (code != Code.OK) {
                    console.log("failed to create account", code, msg, trailers);
                }
            }
        });
    }
    getAccount(){
        let aReq = new Account();
        aReq.setUsername((<HTMLInputElement>document.getElementById("get_username")).value);

        grpc.invoke(AccountService.Get, {
            request: aReq,
            host: this.host,
            onMessage: (message: Account) => {
                this.currAccount = message;
			    this.update();	
            },
            onEnd: (code: Code, msg: string | undefined, trailers: Metadata) => {
                if (code != Code.OK) {
                    console.log("failed to get account", code, msg, trailers);
                    this.currAccount = new Account();
                    this.update();
                }
            }
        });
    }
    login(){
        let aReq = new Account();
        aReq.setUsername((<HTMLInputElement>document.getElementById("login_username")).value);
        aReq.setPassword((<HTMLInputElement>document.getElementById("login_password")).value);

        grpc.invoke(AccountService.Authenticate, {
            request: aReq,
            host: this.host,
            onMessage: (message: AuthenticateAccountResp) => {
                if (message.getLoggedin()){
                    console.log("logged in!");
                } else {
                    console.log("bad username/password");
                }
            },
            onEnd: (code: Code, msg: string | undefined, trailers: Metadata) => {
                if (code != Code.OK) {
                    console.log("failed to login!", code, msg, trailers);
                }
            }
        });
        /*
        grpc.unary(AccountService.Authenticate, {
            request: aReq,
            host: this.host,
			onEnd: res => {
			    const { code, statusMessage, headers, msg, trailers } = res;
			    if (code != Code.OK){
                    console.log("failed to login!", code, msg, trailers);
				} else if (msg.getLoggedin()){
                    console.log("logged in!");
                } else {
                    console.log("bad username/password");
                }
            }
        */
    }
};

var api = new AccountsAPI();
var createAccountButton = document.getElementById('create_account_button');
var getAccountButton = document.getElementById('get_account_button');
var loginButton = document.getElementById('login_button');

createAccountButton.onclick = function(){
	api.create();
}
getAccountButton.onclick = function(){
    api.getAccount();
}
loginButton.onclick = function(){
    api.login();
};

