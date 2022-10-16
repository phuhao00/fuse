# fuse
 
In Game Server ,usually, which have many modules and every module have their logic resolve ,
 
we often bind the module's object with many handlers in each local module ,  then when the 

program start ,the same object's instance include same handler information which we can reuse  in the  runtime .

Although, we may choose do not bind object with many handlers  ,and register these  handlers at local module ,

that's also a good way. 


About fuse ,we can use  like mux ,which also  can resolve the problem

and we can use middlewares like mux , easy to expand.


* 1.reference the idea of mux,but fuse support tcp

* 2.simple,lightweight

* 3.make module independent

* 4.reduce using resource of the runtime 

