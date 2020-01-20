from graphviz import Digraph
import json
import sys

graph = Digraph(comment='automation flow')

filename = sys.argv[1]
with open(filename) as json_file:
    data = json.load(json_file)
    for step in data['steps']:
    	if step.get('action', None) == None:
    		continue

    	graph.node(step['id'], step['id'])
    	if step.get('next', None) != None:
    		if step.get('serverAction', None) != None:
    			graph.edge(step['id'], step['next'], label = 'server:{}'.format(step.get("serverAction")))
    		else:
		    	graph.edge(step['id'], step['next'], label = step.get('action', ""))
    	if step.get('onError', None) != None:		
	    	graph.edge(step['id'], step['onError'], label = "error")

print(graph.source)
graph.render('{}.gv'.format(filename), view=True)


