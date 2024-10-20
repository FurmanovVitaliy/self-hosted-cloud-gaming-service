import Page from "@/components/templates/page";
import authApi from "@/api/auth";
import ResponseError from "@/api/err";
import { eventColorizeInput, hidePassword } from "@/common/input-validations";

import eventBus from "@/common/event-bus";
import constants from "@/common/constants";

//import "@css/login.css";

type loginResponseData = { id: string; username: string; access_token: string };

class LoginPage extends Page {
	private html = (
		<null>
			<section class='auth'>
				<a class='link-button auth__close-button' amepad-focus-element href={constants.routes.home} inner-link animated-upd>
					<span>
						<i class='icon icon--go-close'></i>
					</span>
				</a>
				<div class='auth-container'>
					<form action='' animated-upd>
						<h2>Login</h2>
						<div class='auth-container__input input-box'>
							<label>Email</label>
							<input email type='email' required placeholder='Email'></input>
						</div>
						<div class='auth-container__input input-box'>
							<label>Password</label>
							<div class='input-box__input'>
								<input password type='password' required autocomplete placeholder='Passworld'></input>
								<button type='button'>
									<i class='icon icon--eye-off'></i>
								</button>
							</div>
							<a href='/else'>Forget Password ?</a>
						</div>
						<button class='auth-container__button button button--main' type='submit'>
							Login
						</button>
						<p>
							Don't have an account?{" "}
							<a href={constants.routes.singup} inner-link>
								SingUp
							</a>
						</p>
					</form>
				</div>
			</section>
		</null>
	);
	public title: string = "Authentification | Pixel Cloud";
	//declare main fields
	constructor() {
		super();
		//add content to page
		this.loadTemplate(this.html);
		//local events
		this.addEventListener("input[email]", "input", eventColorizeInput);
		this.addEventListener(".input-box__input button", "click", hidePassword);
		this.addEventListener('button[type="submit"]', "click", this.handleSubmit);
		//local functions
		//connectedCallback
	}
	//! all appeal to html documents as this.container

	private async handleSubmit(e: Event) {
		const emailInput = document.querySelector("input[email]") as HTMLInputElement;
		const passwordInput = document.querySelector("input[password]") as HTMLInputElement;

		e.preventDefault();
		if (!emailInput.value || !passwordInput.value) {
			alert("Please provide both email and password.");
			return;
		}

		const submitButton = document.querySelector('button[type="submit"]') as HTMLButtonElement;
		if (submitButton) {
			submitButton.disabled = true;
			submitButton.textContent = "Logging in...";
		}
		try {
			const responseData = (await authApi.login(emailInput.value, passwordInput.value)) as loginResponseData;
			localStorage.setItem("access_token", responseData.access_token);
			localStorage.setItem("userId", responseData.id);
			localStorage.setItem("username", responseData.username);
			eventBus.notify("nav", constants.routes.home);
		} catch (error) {
			if (error instanceof ResponseError) {
				alert(error.details);
			} else {
				alert("An error occurred. Please try again later.");
			}
		} finally {
			if (submitButton) {
				submitButton.disabled = false;
				submitButton.textContent = "Log in";
			}
		}
	}
	async render() {
		return this.container;
	}
}
export default LoginPage;
