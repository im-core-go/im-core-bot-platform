package grpc

import (
	"context"
	"errors"
	"io"

	"github.com/im-core-go/im-core-bot-platform/internal/logic/chat"
	"github.com/im-core-go/im-core-bot-platform/internal/logic/chat/impls/openai"
	"github.com/im-core-go/im-core-bot-platform/internal/svc"
	chatv1 "github.com/im-core-go/im-core-proto/gen/bot/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatServer struct {
	chatv1.UnimplementedChatServiceServer
	logic chat.Logic
}

func NewChatServer(svcCtx *svc.Context) (*ChatServer, error) {
	logic, err := openai.NewChatLogic(svcCtx)
	if err != nil {
		return nil, err
	}
	return &ChatServer{logic: logic}, nil
}

func (s *ChatServer) PullModels(ctx context.Context, _ *emptypb.Empty) (*chatv1.ModelListResp, error) {
	resp, err := s.logic.PullModules(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "pull models failed: %v", err)
	}
	out := &chatv1.ModelListResp{Data: make([]*chatv1.ModelInfo, 0, len(resp.Data))}
	for _, m := range resp.Data {
		out.Data = append(out.Data, &chatv1.ModelInfo{
			Id:        m.ID,
			CreatedAt: m.CreatedAt,
		})
	}
	return out, nil
}

func (s *ChatServer) CreateConversation(ctx context.Context, req *chatv1.CreateConversationReq) (*chatv1.CreateConversationResp, error) {
	in := &chat.CreateConversationReq{
		Model: req.Model,
		Message: chat.Message{
			Role:        req.GetMessage().GetRole(),
			ContentType: req.GetMessage().GetContentType(),
			Content:     req.GetMessage().GetContent(),
			Meta:        req.GetMessage().GetMeta(),
		},
	}
	resp, err := s.logic.CreateConversation(ctx, in, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create conversation failed: %v", err)
	}
	return &chatv1.CreateConversationResp{
		ConversationId: resp.ConversationID,
		Title:          resp.Title,
		Reply: &chatv1.Message{
			Role:        resp.Reply.Role,
			ContentType: resp.Reply.ContentType,
			Content:     resp.Reply.Content,
			Meta:        resp.Reply.Meta,
		},
	}, nil
}

func (s *ChatServer) ListConversations(ctx context.Context, req *chatv1.ListConversationsReq) (*chatv1.ListConversationsResp, error) {
	in := &chat.ListConversationsReq{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
	}
	resp, err := s.logic.ListConversations(ctx, in, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "list conversations failed: %v", err)
	}
	out := &chatv1.ListConversationsResp{
		Total:    resp.Total,
		Page:     int32(resp.Page),
		PageSize: int32(resp.PageSize),
		Items:    make([]*chatv1.ConversationItem, 0, len(resp.Items)),
	}
	for _, item := range resp.Items {
		out.Items = append(out.Items, &chatv1.ConversationItem{
			ConversationId: item.ConversationID,
			Title:          item.Title,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return out, nil
}

func (s *ChatServer) ListMessages(ctx context.Context, req *chatv1.ListMessagesReq) (*chatv1.ListMessagesResp, error) {
	in := &chat.ListMessagesReq{
		ConversationID: req.GetConversationId(),
		Page:           int(req.GetPage()),
		PageSize:       int(req.GetPageSize()),
	}
	resp, err := s.logic.ListMessages(ctx, in, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "list messages failed: %v", err)
	}
	out := &chatv1.ListMessagesResp{
		Total:    resp.Total,
		Page:     int32(resp.Page),
		PageSize: int32(resp.PageSize),
		Items:    make([]*chatv1.MessageItem, 0, len(resp.Items)),
	}
	for _, item := range resp.Items {
		out.Items = append(out.Items, &chatv1.MessageItem{
			Id:          item.ID,
			Sequence:    item.Sequence,
			Role:        item.Role,
			ContentType: item.ContentType,
			Content:     item.Content,
			Meta:        item.Meta,
			IsSummary:   item.IsSummary,
			CreatedAt:   item.CreatedAt,
		})
	}
	return out, nil
}

func (s *ChatServer) GetConversation(ctx context.Context, req *chatv1.GetConversationReq) (*chatv1.ConversationItem, error) {
	in := &chat.GetConversationReq{ConversationID: req.GetConversationId()}
	resp, err := s.logic.GetConversation(ctx, in, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "get conversation failed: %v", err)
	}
	return &chatv1.ConversationItem{
		ConversationId: resp.ConversationID,
		Title:          resp.Title,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	}, nil
}

func (s *ChatServer) UpdateConversationTitle(ctx context.Context, req *chatv1.UpdateConversationTitleReq) (*emptypb.Empty, error) {
	in := &chat.UpdateConversationTitleReq{
		ConversationID: req.GetConversationId(),
		Title:          req.GetTitle(),
	}
	if err := s.logic.UpdateConversationTitle(ctx, in, req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "update title failed: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ChatServer) DeleteConversation(ctx context.Context, req *chatv1.DeleteConversationReq) (*emptypb.Empty, error) {
	in := &chat.DeleteConversationReq{ConversationID: req.GetConversationId()}
	if err := s.logic.DeleteConversation(ctx, in, req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "delete conversation failed: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ChatServer) ClearMessages(ctx context.Context, req *chatv1.ClearMessagesReq) (*emptypb.Empty, error) {
	in := &chat.ClearMessagesReq{ConversationID: req.GetConversationId()}
	if err := s.logic.ClearMessages(ctx, in, req.GetUserId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "clear messages failed: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ChatServer) Stream(req *chatv1.Completion, srv chatv1.ChatService_StreamServer) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "missing request")
	}
	in := &chat.Completion{
		ConversationID: req.GetConversationId(),
		Model:          req.GetModel(),
		Stream:         req.GetStream(),
		Messages:       make([]chat.Message, 0, len(req.GetMessages())),
	}
	for _, m := range req.GetMessages() {
		in.Messages = append(in.Messages, chat.Message{
			Role:        m.GetRole(),
			ContentType: m.GetContentType(),
			Content:     m.GetContent(),
			Meta:        m.GetMeta(),
		})
	}
	stream, conversationID, err := s.logic.ResponseStream(srv.Context(), in, req.GetUserId())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "stream failed: %v", err)
	}
	defer stream.Close()

	for {
		ev, done, err := stream.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return status.Errorf(codes.Internal, "stream next failed: %v", err)
		}
		if done && ev.ConversationID == "" {
			ev.ConversationID = conversationID
		}
		if sendErr := srv.Send(toProtoStreamEvent(ev)); sendErr != nil {
			return sendErr
		}
		if done {
			return nil
		}
	}
}

func toProtoStreamEvent(ev chat.StreamEvent) *chatv1.StreamEvent {
	return &chatv1.StreamEvent{
		Type:           toProtoStreamEventType(ev.Type),
		Delta:          ev.Delta,
		ConversationId: ev.ConversationID,
		Title:          ev.Title,
	}
}

func toProtoStreamEventType(t chat.StreamEventType) chatv1.StreamEventType {
	switch t {
	case chat.EventTextDelta:
		return chatv1.StreamEventType_STREAM_EVENT_TYPE_TEXT_DELTA
	case chat.EventImage:
		return chatv1.StreamEventType_STREAM_EVENT_TYPE_IMAGE
	case chat.EventDone:
		return chatv1.StreamEventType_STREAM_EVENT_TYPE_DONE
	case chat.EventError:
		return chatv1.StreamEventType_STREAM_EVENT_TYPE_ERROR
	default:
		return chatv1.StreamEventType_STREAM_EVENT_TYPE_UNSPECIFIED
	}
}
