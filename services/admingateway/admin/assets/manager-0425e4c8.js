import{r as w,e as o,c as a,b as t,F as f,d as _,u,k as S,w as E,_ as k,A as b,i as C,o as i,t as r,h as x}from"./index-fd7d3cea.js";const y={class:"mb-4"},$={class:"header cim-header"},A={class:"table cim-table"},B={class:"row g-2 cim-row"},T={class:"form-floating"},F={class:"form-select"},L=["value"],q={__name:"manager",setup(N){let v=x(),d={name:"",count:0},s=w({apps:[],isShowEdit:!1,licenses:[{value:1,label:"0-200"},{value:2,label:"200-500"},{value:3,label:"500-1000"}],app:o.clone(d)});function c(n){s.isShowEdit=n}function h(){let n={name:"",count:""};b.create(n).then(({code:e,msg:l})=>{let p="error",m=l;o.isEqual(e,C.SUCCESS_0.code)&&(p="success",m="创建成功",s.apps.push(n),s.app=o.clone(d)),v.proxy.$toast({icon:p,text:m,duration:4e3})})}function g(){b.getList().then(({data:{items:n}})=>{let e=n.map(l=>(l.created_time=o.formatTime(l.created_time),l.ended_time=o.formatTime(l.ended_time),l.user_count=o.numberWithCommas(l.user_count),l));s.apps=e})}return g(),(n,e)=>(i(),a("div",y,[t("div",$,[e[4]||(e[4]=t("div",{class:"cim-title"},"应用列表",-1)),t("div",{class:"cicon cicon-add cim-button cim-button-bg",onClick:e[0]||(e[0]=l=>c(!0)),onSave:e[1]||(e[1]=l=>h())},"创建应用",32)]),t("table",A,[e[5]||(e[5]=t("thead",null,[t("tr",null,[t("th",{scope:"col"},"应用名称"),t("th",{scope:"col"},"授权个数"),t("th",{scope:"col"},"到期时间"),t("th",{scope:"col"},"创建时间"),t("th",{scope:"col"},"操作")])],-1)),t("tbody",null,[(i(!0),a(f,null,_(u(s).apps,l=>(i(),a("tr",null,[t("td",null,r(l.app_name),1),t("td",null,r(l.user_count),1),t("td",null,r(l.ended_time),1),t("td",null,r(l.created_time),1),t("td",null,[t("a",{class:"btn-link cim-btn-link",type:"button",onClick:e[2]||(e[2]=()=>{})},"查看")])]))),256))])]),S(k,{show:u(s).isShowEdit,title:"创建应用",onHide:e[3]||(e[3]=l=>c(!1))},{default:E(()=>[t("div",B,[e[7]||(e[7]=t("div",{class:"form-floating"},[t("input",{class:"form-control",placeholder:"应用名称"}),t("label",null,"应用名称")],-1)),t("div",T,[t("select",F,[(i(!0),a(f,null,_(u(s).licenses,l=>(i(),a("option",{value:l.value},r(l.label),9,L))),256))]),e[6]||(e[6]=t("label",null,"授权个数",-1))])])]),_:1},8,["show"])]))}};export{q as default};