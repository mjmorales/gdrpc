extends Node

class_name RpcService

var ws: WebSocketClient

func _initialize(webSocketURL: string, webSocketClient: WebSocketClient):
    self.ws = webSocketClient
    self.ws.connect("connection_established", _on_connection_established)
    self.ws.connect("connection_error", _on_connection_error)
    self.ws.connect("data_received", _on_data_received)
    self.ws.connect_to_url(webSocketURL)

func _on_connection_established(protocol):
    print("Connected with protocol: " + protocol)

func _on_connection_error():
    print("Connection failed")

func _on_data_received():
    var data = parse_json(self.ws.get_peer(1).get_packet().get_string_from_utf8())
    print("Received: ", data)

func send_data(data: Dictionary):
     self.ws.get_peer(1).put_packet(to_json(data).to_utf8())