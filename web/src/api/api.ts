import  ResponseError  from "@/api/err";

export abstract class Api {
	constructor(
		protected readonly serverURL: string,
		protected readonly headers: Record<string, string> = {}
	) {}

	protected async request<T>(
		credentials: "include" | "same-origin" | "omit",
		url: string,
		method: string,
		body?: Record<string, any>
	): Promise<T> {
		const options: RequestInit = {
			credentials,
			method,
			headers: {
				"Content-Type": "application/json",
				...this.headers,
			},
		};

		if (body) {
			options.body = JSON.stringify(body);
		}

		try {
			const response = await fetch(`${this.serverURL}${url}`, options);
			const r = response.headers.get('content-type')?.includes('application/json') 
				? await response.json() 
				: null;
		
			if (!response.ok) {
				const errorInfo = {
					message: `Request failed with status code ${response.status}`,
					code: r?.code || '',
					status: response.status,
					details: r?.message || '',
				};
				throw new ResponseError(errorInfo.message, errorInfo.status, errorInfo.details, errorInfo.code);
			}
			
			return r as T;
		
		} catch (error) {
			console.error("Request failed:", error, "URL:", `${this.serverURL}${url}`);
			throw error;  
		}
	}

	protected async get<T>(url: string): Promise<T> {
		return this.request<T>("include", url, "GET");
	}

	protected async post<T>(url: string, body: Record<string, any>): Promise<T> {
		return this.request<T>("same-origin", url, "POST", body);
	}

	protected async put<T>(url: string, body: Record<string, any>): Promise<T> {
		return this.request<T>("same-origin", url, "PUT", body);
	}

	protected async delete<T>(url: string): Promise<T> {
		return this.request<T>("same-origin", url, "DELETE");
	}
}

export default Api;
