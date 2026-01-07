package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	httpmodel "web-chat/api/http_model/chat"
	"web-chat/internal/logic/chat"
	"web-chat/internal/logic/chat/memory"
	"web-chat/internal/svc"
	http2 "web-chat/pkg/http"
	"web-chat/pkg/logger"
	"web-chat/pkg/utils"
)

type logicImpl struct {
	svcCtx  *svc.Context
	utils   *utils.Utils
	urls    *urls
	authKey string
	headers map[string]string
	memory  memory.Manager
}

const (
	titleMaxChars       = 20
	defaultPageSize     = 20
	maxPageSize         = 100
	titleMessageLimit   = 4
)

type completionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionRequest struct {
	Model    string              `json:"model"`
	Messages []completionMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type completionResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewChatLogic(svcCtx *svc.Context) (chat.Logic, error) {
	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		return nil, errors.New("empty openai key")
	}
	baseURL := svcCtx.Config.LLMRequestConf.OpenAI.BaseURL
	if baseURL == "" {
		return nil, errors.New("empty openai base url")
	}
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}
	return &logicImpl{
		svcCtx:  svcCtx,
		utils:   svcCtx.Utils,
		urls:    newURLs(baseURL),
		authKey: "Bearer " + key,
		headers: headers,
		memory: memory.NewManager(
			svcCtx.Dao.ChatDao,
			func() int64 { return svcCtx.Utils.SnowFlake.Generate().Int64() },
			func() string { return svcCtx.Utils.UUID.New() },
		),
	}, nil
}

func (l *logicImpl) ResponseStream(ctx context.Context, req *httpmodel.Completion, userID string) (httpmodel.MessageSteam, string, error) {
	if req.Model == "" {
		return nil, "", errors.New("missing model")
	}
	if len(req.Messages) == 0 {
		return nil, "", errors.New("empty message")
	}

	conversationID, err := l.memory.EnsureConversation(ctx, userID, req.ConversationID)
	if err != nil {
		return nil, "", err
	}
	req.ConversationID = conversationID

	lastInput := req.Messages[len(req.Messages)-1]
	userMsg, err := l.memory.SaveUserMessage(ctx, req.ConversationID, memory.MessageInput{
		Role:        lastInput.Role,
		ContentType: lastInput.ContentType,
		Content:     lastInput.Content,
		Meta:        lastInput.Meta,
	})
	if err != nil {
		return nil, "", err
	}
<<<<<<< Updated upstream
	return newOpenAIResponsesStream(sr), nil
=======

	promptMessages, err := l.memory.BuildPrompt(ctx, req.ConversationID, userMsg, req.Model, func(ctx context.Context, modelName string, messages []memory.PromptMessage) (string, error) {
		return l.doCompletionFromPrompt(ctx, modelName, messages)
	})
	if err != nil {
		return nil, "", err
	}

	sr, err := l.doStreamCompletion(ctx, req.Model, promptMessages)
	if err != nil {
		return nil, "", err
	}

	stream := newOpenAIChatCompletionsStream(sr)
	var streamWithStore httpmodel.MessageSteam
	streamWithStore = newPersistedStream(stream, func(content string) error {
		if err := l.memory.SaveAssistantMessage(ctx, req.ConversationID, content); err != nil {
			return err
		}
		if title, ok := l.generateTitle(req.ConversationID, req.Model); ok {
			l.setStreamTitle(streamWithStore, title)
		}
		return nil
	})
	l.setStreamContext(req.ConversationID, streamWithStore)
	return streamWithStore, req.ConversationID, nil
>>>>>>> Stashed changes
}

func (l *logicImpl) PullModules(ctx context.Context) (*httpmodel.ModelListResp, error) {
	var (
		res = new(httpmodel.ModelListResp)
		err error
	)

	resp, err := l.utils.RequestHandler.DoCommon(
		ctx,
		http.MethodGet,
		l.urls.ModelList,
		nil,
		l.headers,
	)

	if err != nil {
		logger.L().Errorf("pull modules error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (l *logicImpl) CreateConversation(ctx context.Context, req *httpmodel.CreateConversationReq, userID string) (*httpmodel.CreateConversationResp, error) {
	if req.Model == "" {
		return nil, errors.New("missing model")
	}
	if strings.TrimSpace(req.Message.Content) == "" {
		return nil, errors.New("empty message")
	}

	conversationID, err := l.memory.EnsureConversation(ctx, userID, "")
	if err != nil {
		return nil, err
	}

	userMsg, err := l.memory.SaveUserMessage(ctx, conversationID, memory.MessageInput{
		Role:        req.Message.Role,
		ContentType: req.Message.ContentType,
		Content:     req.Message.Content,
		Meta:        req.Message.Meta,
	})
	if err != nil {
		return nil, err
	}

	reply, err := l.doCompletionFromPrompt(ctx, req.Model, []memory.PromptMessage{
		{Role: userMsg.Role, Content: userMsg.Content},
	})
	if err != nil {
		return nil, err
	}

	if err := l.memory.SaveAssistantMessage(ctx, conversationID, reply); err != nil {
		return nil, err
	}

	l.generateTitleAsync(conversationID, req.Model)

	return &httpmodel.CreateConversationResp{
		ConversationID: conversationID,
		Title:          "New",
		Reply: httpmodel.Message{
			Role:        "assistant",
			ContentType: "text",
			Content:     reply,
		},
	}, nil
}

func (l *logicImpl) ListConversations(ctx context.Context, req *httpmodel.ListConversationsReq, userID string) (*httpmodel.ListConversationsResp, error) {
	if userID == "" {
		return nil, errors.New("missing user")
	}
	page, pageSize := normalizePaging(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	items, total, err := l.memory.ListConversations(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	respItems := make([]httpmodel.ConversationItem, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, httpmodel.ConversationItem{
			ConversationID: item.UUID,
			Title:          item.Title,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return &httpmodel.ListConversationsResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    respItems,
	}, nil
}

func (l *logicImpl) ListMessages(ctx context.Context, req *httpmodel.ListMessagesReq, userID string) (*httpmodel.ListMessagesResp, error) {
	if userID == "" {
		return nil, errors.New("missing user")
	}
	if req.ConversationID == "" {
		return nil, errors.New("missing conversation_id")
	}
	conversation, err := l.memory.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if conversation.UserID != userID {
		return nil, errors.New("forbidden")
	}

	page, pageSize := normalizePaging(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	items, total, err := l.memory.ListMessages(ctx, req.ConversationID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	respItems := make([]httpmodel.MessageItem, 0, len(items))
	for _, item := range items {
		meta := ""
		if item.Meta != nil {
			meta = *item.Meta
		}
		respItems = append(respItems, httpmodel.MessageItem{
			ID:          item.ID,
			Sequence:    item.Sequence,
			Role:        item.Role,
			ContentType: item.ContentType,
			Content:     item.Content,
			Meta:        meta,
			IsSummary:   item.IsSummary,
			CreatedAt:   item.CreatedAt,
		})
	}
	return &httpmodel.ListMessagesResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    respItems,
	}, nil
}

func (l *logicImpl) GetConversation(ctx context.Context, req *httpmodel.GetConversationReq, userID string) (*httpmodel.ConversationItem, error) {
	if req.ConversationID == "" {
		return nil, errors.New("missing conversation_id")
	}
	conversation, err := l.memory.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if userID != "" && conversation.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return &httpmodel.ConversationItem{
		ConversationID: conversation.UUID,
		Title:          conversation.Title,
		CreatedAt:      conversation.CreatedAt,
		UpdatedAt:      conversation.UpdatedAt,
	}, nil
}

func (l *logicImpl) UpdateConversationTitle(ctx context.Context, req *httpmodel.UpdateConversationTitleReq, userID string) error {
	if req.ConversationID == "" {
		return errors.New("missing conversation_id")
	}
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("empty title")
	}
	conversation, err := l.memory.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return err
	}
	if userID != "" && conversation.UserID != userID {
		return errors.New("forbidden")
	}
	return l.memory.UpdateConversationTitle(ctx, req.ConversationID, req.Title)
}

func (l *logicImpl) DeleteConversation(ctx context.Context, req *httpmodel.DeleteConversationReq, userID string) error {
	if req.ConversationID == "" {
		return errors.New("missing conversation_id")
	}
	conversation, err := l.memory.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return err
	}
	if userID != "" && conversation.UserID != userID {
		return errors.New("forbidden")
	}
	if err := l.memory.ClearMessages(ctx, req.ConversationID); err != nil {
		return err
	}
	return l.memory.DeleteConversation(ctx, req.ConversationID)
}

func (l *logicImpl) ClearMessages(ctx context.Context, req *httpmodel.ClearMessagesReq, userID string) error {
	if req.ConversationID == "" {
		return errors.New("missing conversation_id")
	}
	conversation, err := l.memory.GetConversation(ctx, req.ConversationID)
	if err != nil {
		return err
	}
	if userID != "" && conversation.UserID != userID {
		return errors.New("forbidden")
	}
	return l.memory.ClearMessages(ctx, req.ConversationID)
}

func (l *logicImpl) doStreamCompletion(ctx context.Context, modelName string, messages []memory.PromptMessage) (*http2.SSEReader, error) {
	prompt := make([]completionMessage, 0, len(messages))
	for _, msg := range messages {
		prompt = append(prompt, completionMessage{Role: msg.Role, Content: msg.Content})
	}
	body, err := json.Marshal(completionRequest{
		Model:    modelName,
		Messages: prompt,
		Stream:   true,
	})
	if err != nil {
		return nil, err
	}
	return l.utils.RequestHandler.DoSSE(ctx, http.MethodPost, l.urls.Completion, bytes.NewReader(body), l.headers)
}

func (l *logicImpl) doCompletion(ctx context.Context, modelName string, messages []completionMessage) (string, error) {
	body, err := json.Marshal(completionRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return "", err
	}
	resp, err := l.utils.RequestHandler.DoCommon(ctx, http.MethodPost, l.urls.Completion, bytes.NewReader(body), l.headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out completionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("empty completion response")
	}
	return out.Choices[0].Message.Content, nil
}

func (l *logicImpl) generateTitleAsync(conversationID, modelName string) {
	go func() {
		ctx := context.Background()
		if title, ok := l.generateTitle(conversationID, modelName); ok {
			_ = l.memory.UpdateConversationTitle(ctx, conversationID, title)
		}
	}()
}

func (l *logicImpl) generateTitle(conversationID, modelName string) (string, bool) {
	ctx := context.Background()
	conversation, err := l.memory.GetConversation(ctx, conversationID)
	if err != nil {
		return "", false
	}
	if conversation.Title != "" && conversation.Title != "New" {
		return "", false
	}
	titleMessages, err := l.memory.BuildTitleMessages(ctx, conversationID, titleMessageLimit)
	if err != nil || len(titleMessages) == 0 {
		return "", false
	}
	prompt := make([]completionMessage, 0, len(titleMessages)+1)
	prompt = append(prompt, completionMessage{
		Role:    "system",
		Content: fmt.Sprintf("Generate a short title (<=%d chars). Return only the title.", titleMaxChars),
	})
	for _, msg := range titleMessages {
		prompt = append(prompt, completionMessage{Role: msg.Role, Content: msg.Content})
	}
	title, err := l.doCompletion(ctx, modelName, prompt)
	if err != nil {
		return "", false
	}
	title = strings.TrimSpace(strings.Trim(title, `"`))
	if title == "" {
		return "", false
	}
	if err := l.memory.UpdateConversationTitle(ctx, conversationID, title); err != nil {
		return "", false
	}
	return title, true
}

func (l *logicImpl) doCompletionFromPrompt(ctx context.Context, modelName string, messages []memory.PromptMessage) (string, error) {
	prompt := make([]completionMessage, 0, len(messages))
	for _, msg := range messages {
		prompt = append(prompt, completionMessage{Role: msg.Role, Content: msg.Content})
	}
	return l.doCompletion(ctx, modelName, prompt)
}

func normalizePaging(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
