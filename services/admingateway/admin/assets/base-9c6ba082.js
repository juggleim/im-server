import{S as v,l as u,r as f,n as _,c as b,a as o,t as a,u as t,p as S,q as n,i as w,A as g,d as r,o as j}from"./index-fb0c9453.js";const x={class:"mb-4 app-base"},A={class:"row jg-ab-row"},I={class:"col-sm-4"},k={class:"form-control-plaintext"},E={class:"row jg-ab-row"},N={class:"col-sm-4"},y={class:"form-control-plaintext"},O={class:"row jg-ab-row"},R={class:"col-sm-4"},T={class:"jg-secret-text"},B={class:"row jg-ab-row"},P={class:"col-sm-2"},K={__name:"base",setup(V){let i=w(),{currentRoute:{_rawValue:{params:{app_key:c}}}}=i;v.get(u.USER_TOKEN);let l=f({appInfo:{restricted_fields:{}},isShowSecret:!1,currentAppStatus:_.ONLINE});function p(){g.getOne({app_key:c}).then(({data:e})=>{let{cur_user_count:s,max_user_count:m}=e;r.formatProps(e,{count:"num",time:"date"}),e.use_percent=Math.floor(s/m)*100,e.n_app_secret="********************",r.extend(l.appInfo,e)})}p();function d(){let e=!l.isShowSecret;r.extend(l,{isShowSecret:e})}return(e,s)=>(j(),b("div",x,[s[8]||(s[8]=o("ul",{class:"nav nav-underline-border ab-underline-border"},[o("li",{class:"nav-item"},[o("a",{class:"nav-link active jgicon jgicon-product"},"基本信息")])],-1)),o("div",A,[s[0]||(s[0]=o("label",{class:"col-sm-1 col-form-label"},"App 名称",-1)),o("div",I,[o("div",k,a(t(l).appInfo.app_name),1)]),s[1]||(s[1]=o("div",{class:"col-sm-4"},null,-1))]),o("div",E,[s[2]||(s[2]=o("label",{class:"col-sm-1 col-form-label"},"App Key",-1)),o("div",N,[o("div",y,a(t(l).appInfo.app_key),1)]),s[3]||(s[3]=o("div",{class:"col-sm-4"},null,-1))]),o("div",O,[s[4]||(s[4]=o("label",{class:"col-sm-1 col-form-label"},"App Secret",-1)),o("div",R,[o("div",{class:S(["form-control-plaintext jg-app-secret",{redfont:t(l).isShowSecret}])},[o("span",T,a(t(l).isShowSecret?t(l).appInfo.app_secret:t(l).appInfo.n_app_secret),1),o("span",{class:"jgicon jgicon-hide jg-secret-btn",onClick:d})],2)]),s[5]||(s[5]=o("div",{class:"col-sm-4"},null,-1))]),s[9]||(s[9]=n('<div class="row jg-ab-row"><label class="col-sm-1 col-form-label">App 到期时间</label><div class="col-sm-4"><div class="form-control-plaintext">2024-05-31 18:30</div></div><div class="col-sm-4"></div></div>',1)),o("div",B,[s[6]||(s[6]=o("label",{class:"col-sm-1 col-form-label"},"授权数量",-1)),o("div",P," 100 / "+a(t(l).appInfo.n_max_user_count),1),s[7]||(s[7]=o("div",{class:"col-sm-4"},null,-1))]),s[10]||(s[10]=n('<div class="row jg-ab-row"><label class="col-sm-1 col-form-label">App 状态</label><div class="col-sm-4"><div class="form-control-plaintext">已上线</div></div><div class="col-sm-4"></div></div>',1))]))}};export{K as default};