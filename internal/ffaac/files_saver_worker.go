package ffaac

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/utilitywarehouse/finance-fulfilment-archive-api-cli/internal/pb/bfaa"
)

type fileSaverWorker struct {
	faaClient bfaa.BillFulfilmentArchiveAPIClient
	fileChan  <-chan string
	errCh     chan<- error
}

func (f *fileSaverWorker) Run(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case fn, ok := <-f.fileChan:
		if ok {
			if err := f.sendFileToArchiveAPI(ctx, fn); err != nil {
				f.errCh <- err
			}
		}
	}
}

func (f *fileSaverWorker) sendFileToArchiveAPI(ctx context.Context, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", fileName)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err, "failed reading bytes for file %s", fileName)
	}

	_, err = f.faaClient.SaveBillFulfilmentArchive(ctx, &bfaa.SaveBillFulfilmentArchiveRequest{
		Id:      fileName,
		Archive: &bfaa.BillFulfilmentArchive{Data: bytes},
	})
	if err != nil {
		return errors.Wrapf(err, "failed calling the fulfilment archive api for file %s", fileName)
	}
	return nil
}