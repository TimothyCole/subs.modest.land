import Vue from 'vue';
import VueRouter from 'vue-router';
import App from './App.vue';

Vue.config.productionTip = false;
Vue.use(VueRouter);

import Header from './components/Header.vue';
import Footer from './components/Footer.vue';

Vue.component("Header", Header);
Vue.component("Footer", Footer);

import HomePage from './pages/HomePage.vue';
import NotFound from './pages/NotFound.vue';
const router = new VueRouter({
	routes: [
		{ name: "Home", path: '/', component: HomePage },
		{ name: "NotFound", path: '*', component: NotFound }
	],
	mode: 'history'
})

new Vue({
	el: '#app',
	router,
	render: h => h(App)
});