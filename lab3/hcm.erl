-module(hcm).
-export([start/0]).

-include_lib("wx/include/wx.hrl").

start() ->
	State = make_window(),
	loop (State).

make_window() ->
  MainWindow = wx:new(),
  Frame = wxFrame:new(MainWindow, -1, "http client", [{size,{800, 600}}]),
  redraw(null, Frame).

redraw(PanelOld, Frame) ->

	case PanelOld of
		null -> ok;
	_ -> wxPanel:destroy(PanelOld)
	end,

	Panel = wxPanel:new(Frame),

	ResultText = wxTextCtrl:new(Panel, 200, [{size, {600, 400}}, {style, ?wxDEFAULT bor ?wxTE_MULTILINE bor ?wxTE_READONLY}]),
	UrlText = wxTextCtrl:new(Panel, 201, [{size, {200, 30}}]),

	RequestBtn = wxButton:new(Panel, 100, [{label, "Request"}]),
	CleanBtn = wxButton:new(Panel, 101, [{label, "Clean"}]),
	SaveBtn = wxButton:new(Panel, 102, [{label, "Save"}]),
	ExitBtn  = wxButton:new(Panel, 103, [{label, "Exit"}]),

	OuterSizer  = wxBoxSizer:new(?wxHORIZONTAL),
	MainSizer   = wxBoxSizer:new(?wxVERTICAL),

	ButtonSizer = wxBoxSizer:new(?wxHORIZONTAL),

	wxSizer:addSpacer(MainSizer, 10),

	wxSizer:add(ButtonSizer, RequestBtn,  []),
	wxSizer:add(ButtonSizer, CleanBtn,  []),
	wxSizer:add(ButtonSizer, SaveBtn,  []),
	wxSizer:add(ButtonSizer, ExitBtn,  []),

	wxSizer:addSpacer(MainSizer, 10),

	wxSizer:add(MainSizer, UrlText, []),
	wxSizer:add(MainSizer, ResultText, []),
	wxSizer:add(MainSizer, ButtonSizer, []),

	wxSizer:add(OuterSizer, MainSizer, []),

	wxPanel:setSizer(Panel, OuterSizer),

	wxFrame:show(Frame),

	wxFrame:connect(Frame, close_window),
	wxPanel:connect(Panel, command_button_clicked),

	{Frame, ResultText, UrlText}.


loop(State) ->
	{Frame, ResultText, UrlText} = State,
	receive

	#wx{event=#wxClose{}} ->
		wxWindow:destroy(Frame),
		ok;

	#wx{id = 100, event=#wxCommand{type = command_button_clicked}} ->
		inets:start(),
		try	{ok, _} = httpc:request(get, {wxTextCtrl:getValue(UrlText), []}, [], [{sync, false}]) of
			_ ->
				wxTextCtrl:appendText(ResultText, "\nRequest sent\n"),
				receive
					{http, {_, {{_, 200, _}, _, Result}}} ->
					self() ! {result, Result},
					wxTextCtrl:appendText(ResultText, "\nAnswer received:")
				after 3000 ->
					error,
					wxTextCtrl:appendText(ResultText, "\nNo answer\n")
				end,
				loop(State)
		catch
			_:_ -> wxTextCtrl:appendText(ResultText, "\nURL format is http://target.web.address\n"),
			loop(State)
		end;

	#wx{id = 101, event=#wxCommand{type = command_button_clicked}} ->
		wxTextCtrl:clear(ResultText),
		loop(State);

	#wx{id = 102, event=#wxCommand{type = command_button_clicked}} ->
		wxTextCtrl:saveFile(ResultText, [{file, "log.txt"}]),
		loop(State);

	#wx{id = 103, event=#wxCommand{type = command_button_clicked} } ->
		wxWindow:destroy(Frame),
		ok;

	{result, Body} ->
		wxTextCtrl:appendText(ResultText, "\n"++binary_to_list(Body)++"\n"),
		loop(State);

	_ ->
		loop(State)

	end.
