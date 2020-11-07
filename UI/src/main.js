import Vue from "vue";
import App from "./App.vue";

Vue.config.productionTip = false;
/* 
  We render the "Virtual DOM" (Vue's version of DOM) of the App component
  onto the element with id="app"
*/
new Vue({
  render: h => h(App)
}).$mount("#app");
