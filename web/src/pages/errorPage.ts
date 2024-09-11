import Page from '@/components/templates/page';

class ErrorPage extends Page {
    static TextObject = {
        Title: 'Error Page',
    };
    constructor() {
        super();
    }
   async render() { 
        const head = document.createElement('h1');
        head.innerText = ErrorPage.TextObject.Title;
        this.container.appendChild(head);
       return this.container;
    
    }
}
export default ErrorPage;