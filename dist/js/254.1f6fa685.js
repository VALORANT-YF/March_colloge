"use strict";(self["webpackChunkvue_jsbk"]=self["webpackChunkvue_jsbk"]||[]).push([[254],{8254:function(e,s,t){t.r(s),t.d(s,{default:function(){return d}});var a=function(){var e=this,s=e._self._c;return s("div",{staticClass:"tabel"},[s("h2",[e._v("搜索：")]),s("el-input",{attrs:{placeholder:"请输入姓名"},on:{input:e.API},model:{value:e.searchname,callback:function(s){e.searchname=s},expression:"searchname"}}),e.isLoading?s("div",{staticClass:"loading-message"},[e._v("数据正在加载中...")]):e._e(),s("el-tabs",{attrs:{type:"border-card"},model:{value:e.activeName,callback:function(s){e.activeName=s},expression:"activeName"}},e._l(e.persons,(function(t,a){return s("el-tab-pane",{key:a,attrs:{label:t.high_dept_name,name:t.high_dept_name}},[s("div",{staticClass:"table"},[s("el-collapse",{on:{change:e.handleChange},model:{value:e.activeNames,callback:function(s){e.activeNames=s},expression:"activeNames"}},e._l(t.TypeNameEnd,(function(t,a){return s("el-collapse-item",{key:a,staticClass:"table2",attrs:{title:t.dept_name,name:a}},[s("div",{staticClass:"jianshu"},[s("el-button",{staticClass:"danger",attrs:{type:"danger"},on:{click:function(s){return e.updateiswrite(t.dept_id,t.is_write_books)}}},[e._v(" "+e._s(e.getyesorno(t.is_write_books))+" ")])],1),e._l(t.TypeName,(function(t,a){return s("div",{key:a,staticClass:"table1"},[s("div",{staticClass:"name"},[e._v(" "+e._s(t.name)+" ")]),s("div",{staticClass:"mobile"},[e._v(" "+e._s(t.mobile)+" "),s("el-button",{staticClass:"primary",attrs:{type:"primary"},on:{click:function(s){return e.toggleAdminStatus(t.user_id,t.is_boss)}}},[e._v(" "+e._s(e.getButtonLabel(t.is_boss))+" ")])],1)])}))],2)})),1)],1)])})),1)],1)},o=[],n=t(9252),i=t(8280),l={name:"tabel",data(){return{persons:[],activeName:"全职部",activeNames:"",searchname:"",isLoading:!0}},methods:{handleChange(e){console.log(" handleChange",e)},async API(){console.log("调用了API");try{this.isLoading=!0;const e=await(0,n.I_)(this.searchname);this.persons=e.data,this.isLoading=!1,this.$emit("data-loaded",this.persons)}catch(e){console.error("请求出错：",e)}},async toggleAdminStatus(e,s){try{console.log("进入线程");const t=e;console.log("userid:",t),console.log("is_boss:",s);const a=s?0:1,o=await(0,n.Tr)(t,a);console.log(o);const{Code:i,Msg:l}=o;console.log("Code:",i),console.log("Msg:",l),"success"===l&&(this.$alert("修改管理员成功","提示",{confirmButtonText:"确定",callback:e=>{this.$message({type:"info",message:"提示:管理员状态已修改"})}}),this.API())}catch(t){console.error("请求出错：",t)}},async updateiswrite(e,s){try{console.log("进入线程2"),console.log("dept_id:",e),console.log("is_write_books:",s);const t=0===s?1:0,a=await(0,n.F3)(e,t);console.log(a);const{Code:o,Msg:i}=a;console.log("Code2:",o),console.log("Msg2:",i),"success"===i&&(this.$alert("修改简书状态成功","提示操作",{confirmButtonText:"确定",callback:e=>{this.$message({type:"info",message:"提示: 操作成功"})}}),this.API())}catch(t){console.error("请求出错：",t)}}},computed:{getButtonLabel(){return e=>!1===e?"设置为管理员":"取消管理员"},getyesorno(){return e=>0===e?"取消部门简书博客":"需要写简书博客"}},created(){i.N.$on("custom-event",this.API)},mounted(){this.API()}},c=l,r=t(1001),u=(0,r.Z)(c,a,o,!1,null,"65bc6522",null),d=u.exports},8280:function(e,s,t){t.d(s,{N:function(){return o}});var a=t(6369);const o=new a["default"]}}]);
//# sourceMappingURL=254.1f6fa685.js.map