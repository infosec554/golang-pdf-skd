package service

import (
	"context"
	"sync"
)

type BatchResult struct {
	Index int
	Data  []byte
	Error error
}

type BatchProcessor struct {
	pdfService PDFService
	maxWorkers int
}

func NewBatchProcessor(pdfService PDFService, maxWorkers int) *BatchProcessor {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}
	return &BatchProcessor{
		pdfService: pdfService,
		maxWorkers: maxWorkers,
	}
}

func (bp *BatchProcessor) CompressBatch(ctx context.Context, inputs [][]byte) []BatchResult {
	results := make([]BatchResult, len(inputs))
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[index] = BatchResult{Index: index, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			output, err := bp.pdfService.Compress().CompressBytes(data)
			results[index] = BatchResult{
				Index: index,
				Data:  output,
				Error: err,
			}
		}(i, input)
	}

	wg.Wait()
	return results
}

func (bp *BatchProcessor) MergeBatch(ctx context.Context, inputSets [][][]byte) []BatchResult {
	results := make([]BatchResult, len(inputSets))
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup

	for i, inputs := range inputSets {
		wg.Add(1)
		go func(index int, data [][]byte) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[index] = BatchResult{Index: index, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			output, err := bp.pdfService.Merge().MergeBytes(data)
			results[index] = BatchResult{
				Index: index,
				Data:  output,
				Error: err,
			}
		}(i, inputs)
	}

	wg.Wait()
	return results
}

func (bp *BatchProcessor) RotateBatch(ctx context.Context, inputs [][]byte, angle int, pages string) []BatchResult {
	results := make([]BatchResult, len(inputs))
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[index] = BatchResult{Index: index, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			output, err := bp.pdfService.Rotate().RotateBytes(data, angle, pages)
			results[index] = BatchResult{
				Index: index,
				Data:  output,
				Error: err,
			}
		}(i, input)
	}

	wg.Wait()
	return results
}

func (bp *BatchProcessor) WatermarkBatch(ctx context.Context, inputs [][]byte, text string, opts *WatermarkOptions) []BatchResult {
	results := make([]BatchResult, len(inputs))
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[index] = BatchResult{Index: index, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			output, err := bp.pdfService.Watermark().AddWatermarkBytes(data, text, opts)
			results[index] = BatchResult{
				Index: index,
				Data:  output,
				Error: err,
			}
		}(i, input)
	}

	wg.Wait()
	return results
}

func (bp *BatchProcessor) ProtectBatch(ctx context.Context, inputs [][]byte, password string) []BatchResult {
	results := make([]BatchResult, len(inputs))
	semaphore := make(chan struct{}, bp.maxWorkers)
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[index] = BatchResult{Index: index, Error: ctx.Err()}
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			output, err := bp.pdfService.Protect().ProtectBytes(data, password)
			results[index] = BatchResult{
				Index: index,
				Data:  output,
				Error: err,
			}
		}(i, input)
	}

	wg.Wait()
	return results
}

type Pipeline struct {
	pdfService PDFService
	operations []PipelineOp
}

type PipelineOp struct {
	Type   string
	Params map[string]interface{}
}

func NewPipeline(pdfService PDFService) *Pipeline {
	return &Pipeline{
		pdfService: pdfService,
		operations: make([]PipelineOp, 0),
	}
}

func (p *Pipeline) Compress() *Pipeline {
	p.operations = append(p.operations, PipelineOp{Type: "compress"})
	return p
}

func (p *Pipeline) Rotate(angle int, pages string) *Pipeline {
	p.operations = append(p.operations, PipelineOp{
		Type: "rotate",
		Params: map[string]interface{}{
			"angle": angle,
			"pages": pages,
		},
	})
	return p
}

func (p *Pipeline) Watermark(text string, opts *WatermarkOptions) *Pipeline {
	p.operations = append(p.operations, PipelineOp{
		Type: "watermark",
		Params: map[string]interface{}{
			"text": text,
			"opts": opts,
		},
	})
	return p
}

func (p *Pipeline) Protect(password string) *Pipeline {
	p.operations = append(p.operations, PipelineOp{
		Type: "protect",
		Params: map[string]interface{}{
			"password": password,
		},
	})
	return p
}

func (p *Pipeline) Execute(input []byte) ([]byte, error) {
	result := input

	for _, op := range p.operations {
		var err error

		switch op.Type {
		case "compress":
			result, err = p.pdfService.Compress().CompressBytes(result)
		case "rotate":
			angle := op.Params["angle"].(int)
			pages := op.Params["pages"].(string)
			result, err = p.pdfService.Rotate().RotateBytes(result, angle, pages)
		case "watermark":
			text := op.Params["text"].(string)
			opts, _ := op.Params["opts"].(*WatermarkOptions)
			result, err = p.pdfService.Watermark().AddWatermarkBytes(result, text, opts)
		case "protect":
			password := op.Params["password"].(string)
			result, err = p.pdfService.Protect().ProtectBytes(result, password)
		}

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *Pipeline) Reset() *Pipeline {
	p.operations = make([]PipelineOp, 0)
	return p
}
